d:
cd D:\Users\small\GoProjects\src\github.com\ethereum\go-ethereum\cmd\geth
geth --rpc --nodiscover --datadir  "./datadir/data0"  --port 30303 --rpcapi "db,eth,net,personal,miner,web3" --rpccorsdomain "*" --networkid 666 --ipcdisable --allow-insecure-unlock console 2>>geth.log