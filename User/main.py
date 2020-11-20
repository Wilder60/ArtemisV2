from User.service.app import App
from .service import app

if __name__ == "__main__":
    app = App()
    app.run()