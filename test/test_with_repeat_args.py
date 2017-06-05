# test function for sharing history storage

from __future__ import print_function
from hyperledger.client import Client

import sys
import time
import timeit

API_URL = 'http://127.0.0.1:7050'
DEPLOY_WAIT = 30
TRAN_WAIT = 2
REPEAT_TIME = 5
CHAINCODE_PATH = "http://gopkg.in/RunshanHu/chaincode-example.v0/example"

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
# * python test.py [API_URL=http://127.0.0.1:7050] [chaincode_name]
if __name__ == '__main__':
    if len(sys.argv) not in [1, 2, 3]:
        print("Usage: python test.py ["
              "API_URL=http://127.0.0.1:7050] [chaincode_name]")
        exit()
    # parse the parameters from command line
    chaincode_name = ""
    if len(sys.argv) >= 2:
        API_URL = sys.argv[1]
    if len(sys.argv) >= 3:
        chaincode_name = sys.argv[2]

    c = Client(base_url=API_URL)

    # write to the file
    f = open('test_4peers.dat', 'w')

    print("Checking cluster at {}".format(API_URL))

    # deploy the chaincode
    if not chaincode_name:
        default_args = ["Data01", "0.5,[add,Joe],[delete,Mike]"]
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

    # check if the initial value is correct
    print(">>>Check the initial value (Data01): ")
    values = query(chaincode_name, ["Data01"], validate=True)
    print(values)
    time.sleep(TRAN_WAIT)

    # query 20 times
    print(">>>Query 20 times, and calculate the time: ")
    duration = timeit.repeat(
        "query(chaincode_name, ['Data01'])",
        number = 20,
        repeat = 3,
        setup = "from __main__ import query, chaincode_name")
    print("time used = {}".format(duration))
    f.write(">>>Query 20 times, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # query 200 time
    print(">>>Query 200 times, and calculate the time: ")
    duration = timeit.repeat(
        "query(chaincode_name, ['Data01'])",
        number = 200,
        repeat = 3,
        setup = "from __main__ import query, chaincode_name")
    print("time used = {}".format(duration))
    f.write(">>>Query 200 times, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # query 2000 times
    print(">>>Query 2000 times, and calculate the time: ")
    duration = timeit.repeat(
        "query(chaincode_name, ['Data01'])",
        number = 2000,
        repeat = 3,
        setup = "from __main__ import query, chaincode_name")
    print("time used = {}".format(duration))
    f.write(">>>Query 2000 times, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # write different values (2 tuples)
    print(">>>Write state to ledger with 2 tuples (once): ")
    args = ["Data02", "0.6, [avg,age], [sum, salary]"]
    duration = timeit.repeat(
        "write(chaincode_name, args, validate = True)",
        number = 1,
        repeat = 3,
        setup = "from __main__ import write, chaincode_name, args")
    print("time used = {}".format(duration))
    f.write(">>>Write to ledger with 2 tuples, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # write different values (2 tuples) for 20 times
    print(">>>Write state to ledger with 2 tuples (20 times): ")
    args = ["Data02", "0.6, [avg,age], [sum, salary]"]
    duration = timeit.repeat(
        "write(chaincode_name, args, validate = True)",
        number = 20,
        repeat = 3,
        setup = "from __main__ import write, chaincode_name, args")
    print("time used = {}".format(duration))
    f.write(">>>Write to ledger with 2 tuples for 20 times, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # write different values (4 tuples)
    print(">>>Write state to ledger with 4 tuples (once): ")
    args = ["Data03", "0.6, [avg,age], [sum,salary], [add,name], [delete, address]"]
    duration = timeit.repeat(
        "write(chaincode_name, args, validate = True)",
        number = 1,
        repeat = 3,
        setup = "from __main__ import write, chaincode_name, args")
    print("time used = {}".format(duration))
    f.write(">>>Write to ledger with 4 tuples, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # write different values (4 tuples) for 20 times
    print(">>>Write state to ledger with 4 tuples (20 times): ")
    args = ["Data03", "0.6, [avg,age], [sum,salary], [add,name], [delete, address]"]
    duration = timeit.repeat(
        "write(chaincode_name, args, validate = True)",
        number = 20,
        repeat = 3,
        setup = "from __main__ import write, chaincode_name, args")
    print("time used = {}".format(duration))
    f.write(">>>Write to ledger with 4 tuples for 20 times, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # write different values (8 tuples)
    print(">>>Write state to ledger with 8 tuples (once): ")
    args = ["Data04", "0.6, [avg,age], [sum,salary], \
                            [add,name], [delete, address], \
                            [avg,salary], [query id], \
                            [delete, age], [query, disease]"]

    duration = timeit.repeat(
        "write(chaincode_name, args, validate = True)",
        number = 1,
        repeat = 3,
        setup = "from __main__ import write, chaincode_name, args")
    print("time used = {}".format(duration))
    f.write(">>>Write to ledger with 8 tuples, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # write different values (8 tuples) for 20 times
    print(">>>Write state to ledger with 8 tuples (20 times): ")
    args = ["Data04", "0.6, [avg,age], [sum,salary], \
                            [add,name], [delete, address], \
                            [avg,salary], [query id], \
                            [delete, age], [query, disease]"]

    duration = timeit.repeat(
        "write(chaincode_name, args, validate = True)",
        number = 20,
        repeat = 3,
        setup = "from __main__ import write, chaincode_name, args")
    print("time used = {}".format(duration))
    f.write(">>>Write to ledger with 8 tuples for 20 times, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)
    # write 20 times
    print(">>>Write to the ledger with 8 tuples for 20 time")
    duration = timeit.repeat(
        "write(chaincode_name, args, validate = True)",
        number = 20,
        repeat = 3,
        setup = "from __main__ import write, chaincode_name, args")
    print("time used = {}".format(duration))
    f.write(">>>Write to ledger with 8 tuples for 20 times, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # write 200 times
    print(">>>Write to the ledger with 8 tuples for 200 times")
    duration = timeit.repeat(
        "write(chaincode_name, args, validate = True)",
        number = 200,
        repeat = 3,
        setup = "from __main__ import write, chaincode_name, args")
    print("time used = {}".format(duration))
    f.write(">>>Write to ledger with 8 tuples for 200 times, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)

    # write 2000 times
    print(">>>Write to the ledger with 8 tuples for 2000 time")
    duration = timeit.repeat(
        "write(chaincode_name, args, validate = True)",
        number = 2000,
        repeat = 3,
        setup = "from __main__ import write, chaincode_name, args")
    print("time used = {}".format(duration))
    f.write(">>>Write to ledger with 8 tuples for 2000 times, time used = {}\n".format(duration))
    time.sleep(TRAN_WAIT)


    # check if the value has been written correctly
    # print(">>>Check the initial value (DataC100): ")
    # values = query(chaincode_name, ["DataC100"], validate=True)
    # print(values)

    f.close()
