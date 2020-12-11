import json
from http import HTTPStatus

import bcrypt
from tornado.web import RequestHandler

from service.db import cloud
from service.model import user

class LoginHandler(RequestHandler):

    def initialize(self, database: cloud.Cloud):
        self.database = database

    def post(self):
        request_user = None        
        try:
            request_user = json.loads(
                self.request.body.decode("UTF-8"),
                object_hook=user.from_payload,
            )
        except KeyError as e:
            self.write_error(HTTPStatus.BAD_REQUEST)
            self.finish()

        stored_user = self.database.get_user(request_user.email)
        if not bcrypt.checkpw(
            bytes(request_user.password), bytes(stored_user.password)):
            self.write_error()

        #generate token here
        self.add_header("token", "Token_Value")
        self.set_status(200)
        self.finish()