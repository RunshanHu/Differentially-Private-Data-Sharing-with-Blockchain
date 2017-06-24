# test function for sharing history storage

from __future__ import print_function
from hyperledger.client import Client

import sys
import time

API_URL = 'http://127.0.0.1:7050'
DEPLOY_WAIT = 20
TRAN_WAIT = 5
CHAINCODE_PATH = "http://gopkg.in/RunshanHu/chaincode-example.v0/fullexample"


# query ledger state
def query(chaincode_name, arg_List, validate=False):

    """
    Query a list of values

    :param chaincode_name: the name fo the chaincode
    :param arg_List: List of arguments.
    :return: a list of values
    """

    result, resp = [], {}
    for arg in arg_List:
        resp = c.chaincode_query(chaincode_name=chaincode_name,
                                 function="query",
                                 args=[arg])
        if not validate:
            continue
        if resp['result']['status'] == 'OK':
            result.append(resp['result']['message'])
            continue
        else:
            print("Error when querying")

    if validate:
        return result


# query match test
def query_match_test(chaincode_name, arg_List, validate=False):

    result, resp = [], {}

    resp = c.chaincode_invoke(chaincode_name=chaincode_name,
                                 function="queryMatchTest",
                                 args=arg_List)

    if resp['result']['status'] == 'OK':
        result.append(resp['result']['message'])
    else:
        print("Error when query match testing")

    if validate:
        return result


def write(chaincode_name, arg_List, validate=False):
    #write to the ledger with different args

    result, resp = [], {}

    resp = c.chaincode_invoke(chaincode_name=chaincode_name,
                                  function="write",
                                  args=arg_List)

    if resp['result']['status'] == 'OK':
        result.append(resp['result']['message'])
    else:
        print("Error when writing")

    if validate:
        return result


# Main  Usage:
# * python test.py [API_URL=http://127.0.0.1:7050] will deploy first
# * python test.py [chaincode_name] [API_URL=http://127.0.0.1:7050]
if __name__ == '__main__':
    if len(sys.argv) not in [1, 2, 3]:
        print("Usage: python test.py ["
              "API_URL=http://127.0.0.1:7050] [chaincode_name]")
        exit()
    # parse the parameters from command line
    chaincode_name = ""
    if len(sys.argv) >= 2:
        chaincode_name = sys.argv[1]
    if len(sys.argv) >= 3:
        API_URL = sys.argv[2]

    c = Client(base_url=API_URL)

    # write result to the file
    f = open('test_4peers.dat', 'w')

    print("Checking cluster at {}".format(API_URL))

    # deploy the chaincode
    if not chaincode_name:
        payload = '''{"budget":0.5,"funType":["sum","avg","max","min"],"results":[-1, -1, -1, -1]}'''

        default_args = ["Data01", payload]
        print(">>>Test: deploy the sharing history storing chaincode"
              "(initial state:{})".format(default_args))

        f.write(">>>Test: deploy the sharing history storing chaincode"
              "(initial state:{})\n".format(default_args))

        res = c.chaincode_deploy(chaincode_path=CHAINCODE_PATH, function="init",
                                 args=default_args)
        chaincode_name = res['result']['message']
        assert res['result']['status'] == 'OK'
        print("Successfully deploy chaincode with returned name = " + chaincode_name)
        print("Wait {}s to make sure deployment is done.".format(DEPLOY_WAIT))

        f.write("Successfully deploy chaincode with returned name = " + chaincode_name + "\n")

        time.sleep(DEPLOY_WAIT)

    # check [0.1, sum]
    payload = '''{"budget":0.1, "funType":"sum"}'''
    print(">>>Query test with payload: {}..".format(payload))
    values = query_match_test(chaincode_name, ["Data01", payload], validate=True)
    print(values)
    # check the result after query match test
    time.sleep(TRAN_WAIT)
    print(">>>Check the result: ")
    values = query(chaincode_name, ["Data01"], validate=True)
    print(values)

    # check [0.1, sum]


    f.close()
