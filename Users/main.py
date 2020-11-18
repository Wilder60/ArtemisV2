import tornado.ioloop
import tornado.web

from routes import users

def make_app() -> tornado.web.Application:
    return tornado.web.Application([
        (r"/users", users.UserHandler, dict(None, None)),
    ])

if __name__ == "__main__":
    app = make_app()
    app.listen(8080)
    tornado.ioloop.IOLoop.current().start()