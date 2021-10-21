#! /bin/bash
#Command examples
export RPC_ADDRESS=localhost:50051
echo $RPC_ADDRESS
numan #shows usage/help 
numan summary 
numan add 353-01-1234568  test.com "test carrier" #create number 
numan view 353-01-1234568 #check it 
numan portin 353-01-1234568 30/1/2021 #set a porting in date
numan allocate 353-01-1234568 55 #allocate to ownerID 55
numan view 353-01-1234568 #check it 
numan list_owner 55 #list by owner
numan list 353-01-1234568 #list by number 
numan list 353-01-12345 #list by partial number 
numan portout 353-01-1234568 30/9/2021 #set a porting out date
numan view 353-01-1234568 #check it (currently port out only shown after deallocation) 
numan deallocate 353-01-1234568 #de-allocate the number 
numan view 353-01-1234568 #check it (currently port out only shown after deallocation) 

