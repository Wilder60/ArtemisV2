import json

import tornado.escape
import tornado.web

from Users.db.postgresql import postgresql

class UserHandler(tornado.web.RequestHandler):

    def initialize(self, database: postgresql, logger):
        self.db = database
        self.log = logger


    def get(self):
        self.db.hello()
        pass

    def post(self):
        body = tornado.escape.json_decode(self.request.body)
        print(body)
        self.write("OK")
    
    def patch(self):
        pass

    def delete(self):
        pass

