import sys, getopt

def main(argv):
    inputfile = ''
    opts, args = getopt.getopt(argv,"i:",["ifile="])
    for opt, arg in opts:
        if opt in ("-i", "--ifile"):
            inputfile = arg
    print ('Input file is ', inputfile)

if __name__ == "__main__":
    main(sys.argv[1:])