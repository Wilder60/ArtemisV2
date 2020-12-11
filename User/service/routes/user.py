from http import HTTPStatus
import json
from service.model.user_info import from_payload

from tornado.web import MissingArgumentError, RequestHandler

from service.db.cloud import Cloud
from service.model import user_info

class UserHandler(RequestHandler):

  def initialize(self, database: Cloud) -> None:
    self.database = database

  def get(self):
    user_id = None
    try:
      user_id = self.get_query_argument("id")    
    except:
      self.set_status(HTTPStatus.BAD_REQUEST)
      self.finish()

    user_data = self.database.get_user_info(user_id)
    data = user_data.to_dict()
    self.finish(data)

  def post(self):
    incoming_data = json.loads(
      self.request.body.decode("UTF-8"),
      user_info.from_payload,
    )

    self.database.insert_user_info()
    self.finish()

  def patch(self):
      incoming_data = json.loads(
        self.request.body.decode("UTF-8"),
        user_info.from_payload,
      )

  def delete(self):
      pass