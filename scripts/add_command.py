import sys
import logging
user = sys.argv[2]
if user != "nuuls":
    logging.error(user+ " cannot add commands")
    sys.exit(1)
file = open("./commands/main.lul", "a")
text = " ".join(sys.argv[3:])

spl = text.split(" ", 2)
logging.error(" ".join(spl))

trigger = spl[1]
reply = spl[2]
file.write("\n" + trigger+ " {\n   text: " + reply + "\n}\n")
file.close()
print("added command " + trigger + " NaM")
