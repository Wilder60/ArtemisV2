import tornado


import tornado.web
from tornado.web import RequestHandler

class UserHandler(RequestHandler):

    def initialize(self) -> None:
        pass

    def get(self, user_id):
        pass

    def post(self):
        pass

    def patch(self):
        pass

    def delete(self):
        pass