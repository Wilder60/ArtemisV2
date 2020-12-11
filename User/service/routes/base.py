import tornado.web
from tornado.web import RequestHandler

class BaseHandler(RequestHandler):
    def prepare(self) -> Optional[Awaitable[None]]:
        pass