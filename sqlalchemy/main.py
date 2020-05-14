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

import logging

from flask import Flask, render_template, request, Response, redirect, url_for

import db_accessor as db 



# Remember - storing secrets in plaintext is potentially unsafe. Consider using
# something like https://cloud.google.com/kms/ to help keep secrets secret.

app = Flask(__name__)

logger = logging.getLogger()

@app.before_first_request
def create_tables():
    db.create_tables()


@app.route("/", methods=["GET"])
def index():
    print("index init")
    try:
        funds = db.get_funds()
    except Exception as e:
        logger.exception(e)
        funds = []
    print(funds)

    return render_template(
        "index.html", funds=funds
    )


@app.route("/fund", methods=["GET", "POST"])
def create_fund():
    # Get the team and time the vote was cast.
    print("Create fund called method: ", request.method)
    #print(request)
    if request.method == "POST":
        try:
            print(request)
            name = request.form["name"]
            description = request.form["description"]

            print(name)
            print(description)
            if name and description:
                db.create_fund(name,description)
        except Exception as e:
            print(repr(e))
            return render_template('NewFund.html', action=repr(e), fund=request.form)

        return redirect(url_for('.index'))
    
    return render_template('NewFund.html', action="Create Fund", fund={})

@app.route("/bonds/<fund_name>", methods=["GET","INSERT","REMOVE"])
def bonds_manager(fund_name):
    print("bonds_manager called method: ", request.method)
	msg = ""
	if request.method == "INSERT":
	    try:
		    print(request)
			parsekey = request.form['bond']
			print(parsekey)
			if parsekey and fund_name:
			    db.insert_bond(fund_name, parsekey)
		except Exception as e:
		    print(repr(e))
			msg = repr(e)
	
	elif request.method == "REMOVE":
		try:
		    print(request)
			parsekey = request.form['bond']
			print(parsekey)
			if parsekey and fund_name:
			    db.remove_bond(fund_name, parsekey)
		except Exception as e:
		    print(repr(e))
			msg = repr(e)
	
	else:
	    try:
		    print(request)
			parsekey = request.form['bond']
			print(parsekey)
			if parsekey and fund_name:
			    db.remove_bond(fund_name, parsekey)
		except Exception as e:
		    print(repr(e))
			msg = repr(e)

	if msg:
		fund.fundName = fund_name
		fund.description = msg
	else:
		fund.fundName = fund_name
		fund.description = db.get_fund_description(fund_name)
		fund.bonds = db.get_bonds_in_fund(fund_name)
		
	return render_template('BondList.html', fund=fund)

if __name__ == "__main__":
    app.run(host="127.0.0.1", port=8080, debug=True)
