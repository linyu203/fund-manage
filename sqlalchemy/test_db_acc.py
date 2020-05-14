import db_accessor as db

def display_all_funds():

    funds = db.get_funds()
    for row in funds:
        print(row)


def add_fund(fund, descript):
    db.create_fund(fund, descript)


def add_bond(fund, bond):
    db.insert_bond(fund, bond)


def del_fund(fund):
    db.remove_fund(fund)

def show_fund_bonds(fund):
    rows = db.get_bonds_in_fund(fund)
    for i in rows:
        print(i)


def init_bonds_funds():
    db.create_tables()
    del_fund("First fund")
    add_fund("First fund", "My first fund")
    add_bond("First fund", "GOOGL US EQUITY")
    add_bond("First fund", "AMAZON US EQUITY")
    add_bond("First fund", "APPLE US EQUITY")
    add_bond("First fund", "FACEBOOK US EQUITY")
    add_bond("First fund", "IBM US EQUITY")
    del_fund("Test fund")
    add_fund("Test fund", "The test fund")
    add_bond("Test fund", "GOOGL US EQUITY")

if __name__=="__main__":
    #init_bonds_funds()
    display_all_funds()
    show_fund_bonds("First fund")
    show_fund_bonds("Test fund")


