# Copyright 2018 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import datetime
import logging
import os

#from flask import Flask, render_template, request, Response
import sqlalchemy


# Remember - storing secrets in plaintext is potentially unsafe. Consider using
# something like https://cloud.google.com/kms/ to help keep secrets secret.
db_user = os.environ.get("DB_USER")
db_pass = os.environ.get("DB_PASS")
db_name = os.environ.get("DB_NAME")
cloud_sql_connection_name = os.environ.get("CLOUD_SQL_CONNECTION_NAME")

#app = Flask(__name__)

#logger = logging.getLogger()

# [START cloud_sql_mysql_sqlalchemy_create]
# The SQLAlchemy engine will help manage interactions, including automatically
# managing a pool of connections to your database
db = sqlalchemy.create_engine(
    # Equivalent URL:
    # mysql+pymysql://<db_user>:<db_pass>@/<db_name>?unix_socket=/cloudsql/<cloud_sql_instance_name>
    sqlalchemy.engine.url.URL(
        drivername="mysql+pymysql",
        username=db_user,
        password=db_pass,
        database=db_name,
        query={"unix_socket": "/cloudsql/{}".format(cloud_sql_connection_name)},
    ),
    # ... Specify additional properties here.
    # [START_EXCLUDE]
    # [START cloud_sql_mysql_sqlalchemy_limit]
    # Pool size is the maximum number of permanent connections to keep.
    pool_size=5,
    # Temporarily exceeds the set pool_size if no connections are available.
    max_overflow=2,
    # The total number of concurrent connections for your application will be
    # a total of pool_size and max_overflow.
    # [END cloud_sql_mysql_sqlalchemy_limit]
    # [START cloud_sql_mysql_sqlalchemy_backoff]
    # SQLAlchemy automatically uses delays between failed connection attempts,
    # but provides no arguments for configuration.
    # [END cloud_sql_mysql_sqlalchemy_backoff]
    # [START cloud_sql_mysql_sqlalchemy_timeout]
    # 'pool_timeout' is the maximum number of seconds to wait when retrieving a
    # new connection from the pool. After the specified amount of time, an
    # exception will be thrown.
    pool_timeout=30,  # 30 seconds
    # [END cloud_sql_mysql_sqlalchemy_timeout]
    # [START cloud_sql_mysql_sqlalchemy_lifetime]
    # 'pool_recycle' is the maximum number of seconds a connection can persist.
    # Connections that live longer than the specified amount of time will be
    # reestablished
    pool_recycle=1800,  # 30 minutes
    # [END cloud_sql_mysql_sqlalchemy_lifetime]
    # [END_EXCLUDE]
)
# [END cloud_sql_mysql_sqlalchemy_create]


def create_tables():
    # Create tables (if they don't already exist)
    with db.connect() as conn:
        conn.execute(
            "CREATE TABLE IF NOT EXISTS funds "
            "( name varchar(31) NOT NULL, "
            "description varchar(101) NOT NULL, creation datetime NOT NULL, PRIMARY KEY (name) );"
        )
        conn.execute(
            "CREATE TABLE IF NOT EXISTS bonds "
            "( fundName varchar(31) NOT NULL, "
            "parsekey varchar(41) NOT NULL , creation datetime NOT NULL);"
        )


def get_funds():
    funds = []
    with db.connect() as conn:
        # Execute the query and fetch all results
        recent_votes = conn.execute(
            """
            select f.name, count(b.parsekey) as count, f.description, f.creation
            from funds as f
            left join bonds as b on b.fundName=f.name
            group by f.name
            """
        ).fetchall()
        # Convert the results into a list of dicts representing votes
        for row in recent_votes:
            funds.append({'Name': row[0], 'Count': row[1], 'Description': row[2], 'Creation': row[3]})

        return funds

def create_fund(name, description):
    stmt = sqlalchemy.text(
        "INSERT INTO funds (name, description, creation)" " VALUES (:name, :description, now())"
    )
        # Using a with statement ensures that the connection is always released
        # back into the pool at the end of statement (even if an error occurs)
    with db.connect() as conn:
        conn.execute(stmt, name=name, description=description)


def insert_bond(fundName, parsekey):
    stmt = sqlalchemy.text(
        "INSERT INTO bonds (fundName, parsekey, creation)" " VALUES (:name, :parsekey, now())"
    )
        # Using a with statement ensures that the connection is always released
        # back into the pool at the end of statement (even if an error occurs)
    with db.connect() as conn:
        conn.execute(stmt, name=fundName, parsekey=parsekey)

def remove_bond(fundName, parsekey):
    stmt = sqlalchemy.text(
        "delete from bonds where fundName=:name and parsekey=:parsekey"
    )
        # Using a with statement ensures that the connection is always released
        # back into the pool at the end of statement (even if an error occurs)
    with db.connect() as conn:
        conn.execute(stmt, name=fundName, parsekey=parsekey)

def remove_fund(name):
    stmt1 = sqlalchemy.text(
        "delete from  bonds where fundName=:name"
    )
    stmt2 = sqlalchemy.text(
        "delete from funds where name=:name"
    )
        # Using a with statement ensures that the connection is always released
        # back into the pool at the end of statement (even if an error occurs)
    with db.connect() as conn:
        conn.execute(stmt1, name=name)
        conn.execute(stmt2, name=name)

def get_bonds_in_fund(fundName):
    bonds = []
    stmt = sqlalchemy.text(
            "select parsekey, creation from bonds where fundName=:fname"
            )
    with db.connect() as conn:
        rows = conn.execute(stmt, fname=fundName).fetchall()
        for row in rows:
            bonds.append({"Ticker": row[0], "Inserted Time": row[1]})

    return bonds

def get_fund_description(fundName):
    desc = ""
    stmt = sqlalchemy.text(
            "select description from funds where name=:fname"
            )
    with db.connect() as conn:
        rows = conn.execute(stmt, fname=fundName).fetchall()
        if rows:
            desc = rows[0][0]
    return desc

