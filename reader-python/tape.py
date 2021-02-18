from readerbase import Reader
import json
import os
from termcolor import colored
import datetime
import sys
import time as t
from pytz import timezone 
import pytz


prevPrice = None

def run(symbol, high, consumer):
  while True:
      for message in consumer:
          item = message[6].decode('utf-8')
          if '{' not in item:
              continue
          items = json.loads(item)
          
          for i in items:
             # Show trades
             if i is not None and "ev" in i and i['ev'] == 'T'  and "sym" in i and i['sym'] == symbol:
                enhanceTape(i, symbol, high)

def enhanceTape(row, symbol, high):
   global prevPrice
   color = None
   side = 'SELL'

   if prevPrice is None:
      prevPrice = row["p"]
      return
   if row["p"] >= prevPrice:
      side = 'BUY'
      
   prevPrice = row["p"]
   price = round(float(row['p']), 2)
   size  = int(row['s'])

   
   time = datetime.datetime.utcfromtimestamp(row['t']/1000.0)
   amount = "{:,}".format(int(price * size))
   amountVal = float(price) * float(size)

   if side == 'BUY':
      color = 'white'
   else:
      color = 'yellow'
   if amountVal > float(high):
      print (colored("%s %5s %7s %7s %11s HIGH" % (str(time)[11:-7], symbol, price, size, amount), color))
   else:
      print (colored("%s %5s %7s %7s %11s" % (str(time)[11:-7], symbol, price, size, amount), color))

if __name__ == '__main__':
    filename = "../conf/polygon.yaml"
    filenameEnv = os.getenv('CONFIGFILE')
    if filenameEnv is not None:
        filename = filenameEnv
    reader = Reader(filename)
    consumer = reader.getConsumer()
    high  = 9999999999
    symbol = 'SPY'
    if len(sys.argv) >= 2:
        symbol = sys.argv[1]
    if len(sys.argv) >=3:
        high= sys.argv[2]
    run(symbol, high, consumer) 
