import sys
import configparser

from service.app import App
from service.db import cloud


def load_config():
  if len(sys.argv) != 2:
    raise FileNotFoundError("Config file not found")

  parser = configparser.ConfigParser()
  parser.read(sys.argv[1], encoding="UTF-8")
  return parser

if __name__ == "__main__":
  config = load_config()
  database = cloud.Cloud(config)
  app = App(database=database, logger=None, config=config)
  app.run()
