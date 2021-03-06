package main

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
  	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/ethclient"
	_"github.com/ethereum/go-ethereum/rpc"
	_"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
)


func simed() {
	var chairmanBalance = big.NewInt(100000000000000000)
	var secondBalance = big.NewInt(100000000000000000)
	var thirdBalance = big.NewInt(100000000000000000)

	chairmanKey, _ := crypto.GenerateKey()
	chairmanAuth := bind.NewKeyedTransactor(chairmanKey)
	chairmanAddr := chairmanAuth.From
	secondKey, _ := crypto.GenerateKey()
	secondAuth := bind.NewKeyedTransactor(secondKey)
	secondAddr:=secondAuth.From
	thirdKey, _ := crypto.GenerateKey();
	thirdAuth := bind.NewKeyedTransactor(thirdKey)
	thirdAddr:=thirdAuth.From

	alloc:=core.GenesisAlloc{
		chairmanAddr: {Balance:chairmanBalance},
		secondAddr: {Balance:secondBalance},
		thirdAddr: {Balance:thirdBalance}}
	fmt.Printf("Created simulated backend with chairman %v\n", chairmanAuth.From.Hex())
	conn := backends.NewSimulatedBackend(alloc)


	var context = context.Background()
	balance, err := conn.BalanceAt(context, chairmanAddr, nil)

	if (err != nil) {
		fmt.Printf("chairman balance %v", err)
	}

	fmt.Printf("chairman balance : %v\n", balance)


	var array1  = [32]byte{0:'a',1:'a'}
	var array2  = [32]byte{0:'b',1:'b'}
	var array3  = [32]byte{0:'c',1:'c'}
	var array4  = [32]byte{0:'d',1:'d'}
	var proposalNames =[][32]byte{array1,array2,array3,array4}

	contractAddr, transaction, ballot, err := DeployBallot(chairmanAuth, conn, proposalNames)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v\n", err)
	}
	fmt.Printf("contract addr:%v\n", contractAddr.Hex())
	fmt.Printf("Transaction hash: %v\n", transaction.Hash().Hex());
	// Print the current (non existent) and pending name of the contract
	addr, _ := ballot.Chairperson(nil)
	fmt.Printf("Pre-mining Chairperson addr:%v\n", addr.Hex())

	addr, _ = ballot.Chairperson(&bind.CallOpts{Pending: true})
	fmt.Printf("Pre-mining Chairperson pending addr:%v\n", addr.Hex())

	// Commit all pending transactions in the simulator and print the names again
	conn.Commit()
	addr, _ = ballot.Chairperson(nil)
	fmt.Printf("Post-mining Chairperson addr:%v\n", addr.Hex())

	addr, _ = ballot.Chairperson(&bind.CallOpts{Pending: true})
	fmt.Printf("Post-mining Chairperson pending addr:%v\n", addr.Hex())

	printVoters(ballot,chairmanAddr)
	printVoters(ballot,secondAddr)
	printVoters(ballot,thirdAddr)
	ballot.GiveRightToVote(chairmanAuth,secondAddr)
	ballot.GiveRightToVote(chairmanAuth,thirdAddr)
	ballot.Vote(chairmanAuth,big.NewInt(int64(0)))
	// ballot.Vote(secondAuth,big.NewInt(int64(1)))
	conn.Commit()
	// ballot.Vote(secondAuth,big.NewInt(int64(2)))
	ballot.Vote(thirdAuth,big.NewInt(int64(3)))
	//设置third由second代理投票，分别测试second投票前和投票后的情况
	// ballot.Delegate(thirdAuth,secondAddr)
	//如果second已经投票，则此函数会将third的票累积到second投的人身上，若second未投票，则second.weight+=third。weight
	// ballot.Vote(secondAuth,big.NewInt(int64(1)))
	conn.Commit()
	printVoters(ballot,chairmanAddr)
	// printVoters(ballot,secondAddr)
	// printVoters(ballot,thirdAddr)
	printProposals(ballot)
	winnner,_:=ballot.WinnerName(nil)
	fmt.Println(winnner)
	testaddr,_:=ballot.Testmsg(nil)
	fmt.Printf("testaddr addr:%v\n", testaddr.Hex())
	testdata,_:=ballot.Testdata(nil)
	fmt.Printf("testdata :%v\n", testdata)
	testvalue,_:=ballot.Testvalue(nil)
	fmt.Printf("testvalue :%v\n", testvalue)
	testGasPrice,_:=ballot.TestGasPrice(nil)
	fmt.Printf("testGasPrice:%v\n", testGasPrice)
	testgas,_:=ballot.Testgas(nil)
	fmt.Printf("testGas:%v\n", testgas)
	testorigin,_:=ballot.Testtxorigin(nil)
	fmt.Printf("transaction origin add:%v\n",testorigin.Hex())
}
func main() {
	// testrpc()
	simed()
}
func testrpc() {//conn to testrpc method
	chairmanAddr:=common.HexToAddress("4611f13ccafb5acaf2c1909d26f5dbb7d4d774a0")
	secondAddr:=common.HexToAddress("6a1ae1d9f9ade400154c97d61613e4b0164d5710")
	thirdAddr:=common.HexToAddress("2e1803b56353a2e0f318610cb7f1517d52c247b0")
	chairmanAuth:=&bind.TransactOpts{From:chairmanAddr}
	secondAuth:=&bind.TransactOpts{From:secondAddr}
	thirdAuth:=&bind.TransactOpts{From:thirdAddr}

	conn, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	var context = context.Background()
	balance, err := conn.BalanceAt(context, chairmanAddr, nil)

	fmt.Println("//////////////////////////////////////////////////////")
	if (err != nil) {
		fmt.Printf("chairman balance %v\n", err)
	}

	fmt.Printf("chairman balance : %v\n", balance)


	var array1  = [32]byte{0:'a',1:'a'}
	var array2  = [32]byte{0:'b',1:'b'}
	var array3  = [32]byte{0:'c',1:'c'}
	var array4  = [32]byte{0:'d',1:'d'}
	var proposalNames =[][32]byte{array1,array2,array3,array4}

	contractAddr, transaction, ballot, err := DeployBallot(chairmanAuth, conn, proposalNames)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v\n", err)
	}
	fmt.Printf("contract addr:%v\n", contractAddr.Hex())
	fmt.Printf("Transaction hash: %v\n", transaction.Hash().Hex());
	// Print the current (non existent) and pending name of the contract
	addr, _ := ballot.Chairperson(nil)
	fmt.Printf("Pre-mining Chairperson addr:%v\n", addr.Hex())

	addr, _ = ballot.Chairperson(&bind.CallOpts{Pending: true})
	fmt.Printf("Pre-mining Chairperson pending addr:%v\n", addr.Hex())

	addr, _ = ballot.Chairperson(nil)
	fmt.Printf("Post-mining Chairperson addr:%v\n", addr.Hex())

	addr, _ = ballot.Chairperson(&bind.CallOpts{Pending: true})
	fmt.Printf("Post-mining Chairperson pending addr:%v\n", addr.Hex())

	printVoters(ballot,chairmanAddr)
	printVoters(ballot,secondAddr)
	printVoters(ballot,thirdAddr)
	ballot.GiveRightToVote(chairmanAuth,secondAddr)
	ballot.GiveRightToVote(chairmanAuth,thirdAddr)
	ballot.Vote(chairmanAuth,big.NewInt(int64(0)))
	// ballot.Vote(secondAuth,big.NewInt(int64(1)))
	
	// ballot.Vote(secondAuth,big.NewInt(int64(2)))
	ballot.Vote(thirdAuth,big.NewInt(int64(3)))
	//设置third由second代理投票，分别测试second投票前和投票后的情况
	ballot.Delegate(thirdAuth,secondAddr)
	//如果second已经投票，则此函数会将third的票累积到second投的人身上，若second未投票，则second.weight+=third。weight
	ballot.Vote(secondAuth,big.NewInt(int64(1)))
	
	printVoters(ballot,chairmanAddr)
	printVoters(ballot,secondAddr)
	printVoters(ballot,thirdAddr)
	printProposals(ballot)
	winnner,_:=ballot.WinnerName(nil)
	fmt.Println(winnner)
	testaddr,_:=ballot.Testmsg(nil)
	fmt.Printf("testaddr addr:%v\n", testaddr.Hex())
}
func printProposals(ballot *Ballot) {
	fmt.Println("///////////////////////////////////////")
	for i:=0;i<4;i++{
			proposals, err :=ballot.Proposals(nil,big.NewInt(int64(i)))
			// fmt.Printf("proposals:%s,%d\n", proposals.Name,proposals.VoteCount)
			if err!=nil{
					fmt.Println("error:%v",err)
				}else{
					fmt.Printf("%+v\n", proposals)	
				}
		}
	fmt.Println("///////////////////////////////////////")
}
func printVoters(ballot *Ballot,addr common.Address) {
	fmt.Println("///////////////////////////////////////")
	
	voters, err :=ballot.Voters(nil,addr)
	if err!=nil{
			fmt.Println("error:%v",err)
		}else{
			fmt.Printf("%+v\n", voters)	
		}
}

func rpc() {
	account0Transactor, err := bind.NewTransactor(strings.NewReader(key), "123456")
	if err != nil {
	    log.Fatalf("Failed to create authorized transactor: %v", err)
	  }
	account1Transactor, err := bind.NewTransactor(strings.NewReader(key1), "123456")
	if err != nil {
	    log.Fatalf("Failed to create authorized transactor: %v", err)
	  }

	conn, err := ethclient.Dial("http://localhost:8545")
		if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		}


	var context = context.Background()
	balance, err := conn.BalanceAt(context, common.HexToAddress("58185e1874446969fff65372b404af9a51842156"), nil)

	if (err != nil) {
		fmt.Printf("%v", err)
	}

	fmt.Printf("Balance : %v\n", balance)


	var array1  = [32]byte{0:'a',1:'a'}
	var array2  = [32]byte{0:'b',1:'b'}
	var array3  = [32]byte{0:'c',1:'c'}
	var array4  = [32]byte{0:'d',1:'d'}
	var proposalNames =[][32]byte{array1,array2,array3,array4}

	contractAddr, transaction, ballot, err := DeployBallot(account0Transactor, conn, proposalNames)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v\n", err)
	}
	fmt.Printf("contract addr:%v\n", contractAddr.Hex())
	fmt.Printf("Transaction hash: %v\n", transaction.Hash().Hex());
	// Print the current (non existent) and pending name of the contract
	addr, _ := ballot.Chairperson(nil)
	fmt.Printf("Pre-mining Chairperson addr:%v\n", addr.Hex())

	addr, _ = ballot.Chairperson(&bind.CallOpts{Pending: true})
	fmt.Printf("Pre-mining Chairperson pending addr:%v\n", addr.Hex())

	// Commit all pending transactions in the simulator and print the names again
	//sim.Commit()
	time.Sleep(8000 * time.Millisecond)
	addr, _ = ballot.Chairperson(nil)
	fmt.Printf("Post-mining Chairperson addr:%v\n", addr.Hex())

	addr, _ = ballot.Chairperson(&bind.CallOpts{Pending: true})
	fmt.Printf("Post-mining Chairperson pending addr:%v\n", addr.Hex())

	firstAccount:=common.HexToAddress("390bae9e7e9684a09b1aa73590eee3e78add44a0")
	secondAccount:=common.HexToAddress("58185e1874446969fff65372b404af9a51842156")
	printVoters(ballot,firstAccount)
	printVoters(ballot,secondAccount)
	ballot.GiveRightToVote(account0Transactor,secondAccount)
	time.Sleep(4000 * time.Millisecond)
	ballot.Vote(account0Transactor,big.NewInt(int64(0)))
	time.Sleep(4000 * time.Millisecond)
	ballot.Vote(account1Transactor,big.NewInt(int64(1)))
	time.Sleep(4000 * time.Millisecond)
	ballot.Vote(account0Transactor,big.NewInt(int64(2)))
	time.Sleep(4000 * time.Millisecond)
	ballot.Vote(account0Transactor,big.NewInt(int64(3)))
	time.Sleep(8000 * time.Millisecond)
	printVoters(ballot,firstAccount)
	printVoters(ballot,secondAccount)
	printProposals(ballot)
}