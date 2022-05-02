import time
import traceback
from dataclasses import dataclass
import sqlite3
import argparse
from abc import ABC, abstractmethod
from typing import List, Any

import pandas as pd
import surprise
from aiohttp.web_exceptions import HTTPException
from surprise import KNNBaseline
from surprise import Reader, Dataset
import joblib
from aiohttp import web


class ModelTrainer(ABC):
    model: Any

    @abstractmethod
    def train_model(self):
        raise NotImplemented

    def save_model(self, filename):
        if self.model is None:
            raise ValueError("attempt to save None model.")
        ModelPersistence.save(self.model, filename)


class SurpriseModelTrainer(ModelTrainer):
    def __init__(self, db, data_count=10000, algo=KNNBaseline):
        self._db = db
        self._data_count = data_count
        self._algo = algo

        self._conn = sqlite3.connect(self._db)
        self._pt = None

        self.model: surprise.AlgoBase = None

        self._init_data()

    def _init_data(self):
        self._pt = pd.read_sql(f"SELECT * FROM playlist_tracks LIMIT {self._data_count}",
                               self._conn)

        # 播放列表 -> 用户, 曲目 -> 电影
        self._pt = self._pt.rename(columns={"playlist_id": "userID", "track_id": "itemID"})

        # 歌在播放列表里，就是用户给歌打了一分
        self._pt = self._pt.join(pd.Series([1] * len(self._pt), name="rating"))

        # 洗牌
        self._pt = self._pt.sample(frac=1)

        # surprise: custom dataset
        reader = Reader(rating_scale=(0, 1))
        self._train_data = Dataset.load_from_df(
            self._pt[['userID', 'itemID', 'rating']],
            reader)
        self._train_set = self._train_data.build_full_trainset()

    def train_model(self):
        sim_options = {
            'user_based': False  # compute  similarities between items
        }

        # 训练
        self.model = self._algo(sim_options=sim_options)
        self.model.fit(self._train_set)


train_algos = {'KNNBaseline': KNNBaseline}


class ModelPersistence(object):
    """Model persistence with joblib
    """

    @staticmethod
    def save(model, filename, compress=True):
        """Save model into filename

        :param model: model to save
        :param filename: pah to saving file
        :param compress: Optional compression level for the data. 0 or False is
            no compression. Higher value means more compression, but also slower
            read and write times. Using a value of 3 is often a good compromise.
            If compress is True, the compression level used is 3.
        :return: The list of file names in which the data is stored.
        """
        names = joblib.dump(model, filename, compress=compress)
        print("saved:", names)
        return names

    @staticmethod
    def load(filename):
        """load a saved model

        :param filename: the saved file
        :return: model
        """
        return joblib.load(filename)


@dataclass
class NextSongSeed(object):
    track_name: str
    artists: List[str]


@dataclass
class NextSong(object):
    track_id: str
    track_name: str
    artists: List[str]
    album_cover: str


class TrackMatcher(ABC):
    def __init__(self, tid: str, track_name: str, track_artists: List[str]):
        """tid matches track_name + track_artists
        """
        self.tid = tid
        self.track_name = track_name
        self.track_artists = track_artists

    @abstractmethod
    def is_matched(self, *args, **kwargs) -> bool:
        raise NotImplemented

    def __call__(self, *args, **kwargs):
        return self.is_matched(*args, **kwargs)


class AlwaysMatcher(TrackMatcher):
    """Always returns True. Not recommended.
    """

    def is_matched(self, *args, **kwargs) -> bool:
        return True


class NextSongRecommender(ABC):
    @abstractmethod
    def recommend_next_song(self, seed: NextSongSeed, k=5, shift=None) -> List[NextSong]:
        raise NotImplemented


def sql_sanitize(s: str) -> str:
    """sanitize database inputs texts to prevent SQL injections

    :param s: string to input to sql
    :return: sanitized string
    """
    cs = f'{s}'
    for c in cs:
        if not c.isalnum() and not c.isalpha():
            s = s.replace(c, ' ')
    # punc = r"""!()-[]{};:'"\,<>./?@#$%^&*_~+="""
    # for p in punc:
    #     s = s.replace(p, ' ')
    return s


class SurpriseRecommender(NextSongRecommender):
    def __init__(self, db, model_filename):
        self._db = db
        self._conn = sqlite3.connect(self._db)

        self.model = ModelPersistence.load(filename=model_filename)

    def find_artists(self, track_id: str) -> list:
        c = self._conn.cursor()  # get artists
        c.execute(
            "select a.name from artists a inner join track_artists ta on a.id = ta.artist_id where ta.track_id=?",
            (track_id,))
        a = [a[0] for a in c.fetchall()]
        c.close()
        return a

    def find_track_id(self, name: str, artists: List[str], limit=50, matcher=AlwaysMatcher) -> str:
        """find a track in db.
        Returns the found id, or else raises a ValueError

        :param name: track name
        :param artists: track artists list
        :return: id
        """
        c = self._conn.cursor()
        # 用 FTS 优化: https://www.sqlite.org/fts5.html
        #   CREATE VIRTUAL TABLE tracks_name_fts USING fts5(name, id);
        #   INSERT INTO tracks_name_fts SELECT name, id FROM tracks;
        name = sql_sanitize(name)
        c.execute(f"SELECT id FROM tracks_name_fts WHERE name MATCH '{name}' LIMIT {limit}")
        name_matched_ids = c.fetchall()  # [(id, ), ...]

        for res in name_matched_ids:
            tid = res[0]
            if matcher(tid, name, artists)():
                return tid
        raise ValueError("track not found")

    def find_sim(self, track_id, k=5, shift=None) -> list:
        """ 找和 track_id 曲目最相近的 k 首歌

        :param k: 找到 track_id 的 k 个近邻
        :param shift: 偏移 k 个近邻，避免反复推荐同样的几个东西：
            从模型获取 k+shift 个近邻，然后输出后 shift 个结果（丢弃前 shift 个近邻）。
            default shift=None: `shift = k / 3`
        :return: list of tracks [{id, name, artists, image}, ...] : len=(k+1), [0] 是输入的 track_id
        """

        shift = k // 3 if shift is None else int(shift)

        sim = self.model.get_neighbors(iid=self.model.trainset.to_inner_iid(track_id), k=k + shift)
        sim = sim[shift:]

        c = self._conn.cursor()

        track_ids = [track_id] + list(
            map(self.model.trainset.to_raw_iid, sim))
        # print(track_ids)

        tracks = []
        for tid in track_ids:
            c.execute(f"SELECT * FROM tracks WHERE id = '{tid}'")
            tk = c.fetchall()
            if len(tk) < 1:
                print("fetch track: len(tk) < 1:", tk)
                continue
            tk = tk[0]
            tracks.append(tk + (self.find_artists(tid),))
        c.close()
        # disc_number|duration|explicit|endpoint|id|name|preview_url|track_number|uri|type|image_url | artists
        return [{"id": r[4], "name": r[5], "artists": r[-1], "image": r[-2]} for r in tracks]

    def recommend_next_song(self, seed: NextSongSeed, k=5, shift=None) -> List[NextSong]:
        seed_id = self.find_track_id(seed.track_name, seed.artists)
        recommended = self.find_sim(seed_id, k, shift)
        return list(
            map(lambda x: NextSong(x["id"], x["name"], x["artists"], x["image"]),
                recommended)
        )


class Server:
    def __init__(self, recommender: NextSongRecommender):
        self.recommender = recommender

    async def recommend_handler(self, request: web.Request):
        track_name = request.query.get("track_name") or ""
        artists = request.get("artists") or ""
        k = int(request.query.get("k") or 5)
        shift = int(request.query.get("shift") or -1)
        shift = shift if shift > 0 else None
        # shift = shift if shift < k else k

        if track_name == artists == "":
            raise web.HTTPBadRequest(text="seed track_name and artists expected")

        try:
            recommended = self.recommender.recommend_next_song(
                seed=NextSongSeed(track_name, artists), k=k, shift=shift)
            return web.json_response(list(map(lambda r: r.__dict__, recommended)))
        except ValueError as e:
            not_founds = ['not part of the trainset', 'track not found']
            if any(n in str(e) for n in not_founds):
                raise web.HTTPNotFound(text=str(e))
            print(f'[ERROR] {time.ctime(time.time())} Server got an unknown ValueError:', e)
            traceback.print_exc()
            raise web.HTTPInternalServerError(text=str(e))


# region cors


CORS_HEADERS = {
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': '*',
    'Access-Control-Allow-Headers': '*',
    'Access-Control-Allow-Credentials': 'true',
}


@web.middleware
async def cors_middleware(request, handler):
    """用来解决 cors
    `app = web.Application(middlewares=[cors_middleware])`
    """
    # if request.method == 'OPTIONS':
    #     response = web.Response(text="")
    # else:
    try:
        response = await handler(request)

        for k, v in CORS_HEADERS.items():
            response.headers[k] = v

        return response
    except HTTPException as e:
        for k, v in CORS_HEADERS.items():
            e.headers[k] = v
        raise e


async def empty_handler(request):
    """给每个 route 配上一个对应的 options empty_handler 即可解决 cors 问题:
    `web.options('...', empty_handler)`
    """
    return web.Response()


# endregion cors

# region logger

def log(level, request, response):
    print(f"[{level}] {time.ctime(time.time())} {request.method} {request.rel_url} -> {response.status}")  # log


@web.middleware
async def log_middleware(request, handler):
    try:
        response = await handler(request)
        log('INFO', request, response)
        return response
    except HTTPException as e:
        log('WARN', request, e)
        raise e


# endregion logger


# region cli

def run_service(db: str, model_file: str, host: str, port: int):
    recommender = SurpriseRecommender(db, model_file)
    server = Server(recommender)

    app = web.Application(middlewares=[log_middleware, cors_middleware])
    app.add_routes([
        web.get('/next-song', server.recommend_handler),
        # cors
        web.options('/next-song', empty_handler),
        web.options('/', empty_handler),
    ])

    web.run_app(app, host=host, port=port)


def run_train(db: str, model_file: str, data_count: int = 10000, algo: str = 'KNNBaseline'):
    assert algo in train_algos, 'unknown algo'

    trainer = SurpriseModelTrainer(db, data_count=data_count, algo=train_algos[algo])

    trainer.train_model()

    trainer.save_model(model_file)


if __name__ == '__main__':
    parser = argparse.ArgumentParser("murecom-intro", description="a next-song recommender")
    subparsers = parser.add_subparsers(help="sub commands")

    parser_service = subparsers.add_parser("service")
    parser_service.add_argument("--db", type=str, help="path to SQLite database", required=True)
    parser_service.add_argument("--model", type=str, help="path to save trained model", required=True)
    parser_service.add_argument("--host", type=str, help="server host", default="localhost")
    parser_service.add_argument("--port", type=int, help="listen port", default=8080)

    parser_service.set_defaults(func=lambda args: run_service(args.db, args.model, args.host, args.port))

    parser_train = subparsers.add_parser("train")
    parser_train.add_argument("--db", type=str, help="path to SQLite database", required=True)
    parser_train.add_argument("--model", type=str, help="path to save trained model", required=True)
    parser_train.add_argument("--data-count", type=int, help="data count in the train-set", default=10000)
    parser_train.add_argument("--algo", type=str, help="algorithm to build the model", default="KNNBaseline")
    parser_train.set_defaults(func=lambda args: run_train(args.db, args.model, args.data_count, args.algo))

    args = parser.parse_args()
    try:
        args.func(args)
    except AttributeError:  # 'Namespace' object has no attribute 'func'
        parser.print_usage()

# endregion cli
