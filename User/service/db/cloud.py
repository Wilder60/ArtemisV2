from contextlib import contextmanager

import sqlalchemy
from sqlalchemy.engine.url import URL
from sqlalchemy.orm import sessionmaker
from sqlalchemy.orm import session
from sqlalchemy.orm.exc import NoResultFound, MultipleResultsFound

from service.db import driver
from service.model import login
from service.model import user
from service.model import schema

class Cloud(driver.Base):
    def __init__(self, config) -> None:
        engine = sqlalchemy.create_engine("sqlite:///:memory:", echo=True)

            # URL(
            #     drivername="postgresql+pg8000",
            #     username=config.db.cloud.db_user,
            #     password=config.db.cloud.db_pass,
            #     host=config.db.cloud.db_hostname,
            #     port=config.db.cloud.db_port,
            #     database=config.db.cloud.db_name,
            # ),
        # )

        schema.Base.metadata.create_all(engine, checkfirst=True)
        self.DBSession = sessionmaker(bind = engine)

    @contextmanager
    def session_scope(self) -> session.Session:
      """handles creating the session, and commiting the session

      Creates a session and yields in to the calling code.  If the calling 
      code,successfully completes the session is commit, if an exception is 
      thrown, the transaction is rolledback and the exception is propogated. 
      For example:

        with self.session_scope() as session:
          session.query(User).filter_by(email="foo")

      """
      session = self.DBSession()

      try:
        yield session
        session.commit()

      except:
        session.rollback()

      finally:
        session.close()                

    def insert_user(self, new_user: user.Info) -> None:
      """inserts the user info in the appropiate dbs

      Args:
          new_user: user object from incoming request

      algo:
          get count of login table

      """
      with self.session_scope() as session:
        user_exists = session.query(user.Info).filter_by(email=new_user.email).count()
        if user_exists:
          return

        

        session.add(user)

    def get_user(self, user_email: str) -> user.Info:
        with self.session_scope() as session:
            try:
                db_user = session.query(user.Info).filter_by(email=user_email).one()
                return db_user
            except NoResultFound or MultipleResultsFound as err:
                # if we catch a MultipleResultsFound exception that should be 
                # logged because it's really bad
                return None

    def update_user(self, data: user) -> None:
        session = self.DBSession()
        pass

    def delete_user(self, data: user) -> None:
        pass

    def get_user_info(self, user: str) -> user.Info:
        pass

    def update_user_info(self, data: user.Info) -> None:
        pass

    def delete_user_info(self, data: user.Info) -> None:
        pass