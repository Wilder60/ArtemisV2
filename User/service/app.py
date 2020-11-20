import tornado.web
from tornado.web import Application

class App():
    def __init__(self, database, logger, config) -> None:
        self.database = database
        self.logger = logger
        self.config = config

        self.app = Application([
            (r"/login", None),
            (r"/user", None),
        ]) 

    def run(self):
        self.app.listen(8080)
        tornado.ioloop.IOLoop.current().start()