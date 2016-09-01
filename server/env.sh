if [ -z "$GOPATH" ]; then
	export $GOPATH=$(pwd)
else
	export GOPATH=$GOPATH:$(pwd)
fi

BASE=$(pwd)
ROOT=`dirname $BASE`
export GOPATH=$GOPATH:$ROOT/third
