package main

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
  	_"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/ethclient"
	_"github.com/ethereum/go-ethereum/rpc"
	_"github.com/ethereum/go-ethereum/eth"
	_"github.com/ethereum/go-ethereum/core"
	_"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
)


func main() {
	
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
	time.Sleep(4000 * time.Millisecond)
	addr, _ = ballot.Chairperson(nil)
	fmt.Printf("Post-mining Chairperson addr:%v\n", addr.Hex())

	addr, _ = ballot.Chairperson(&bind.CallOpts{Pending: true})
	fmt.Printf("Post-mining Chairperson pending addr:%v\n", addr.Hex())

	firstAccount:=common.HexToAddress("390bae9e7e9684a09b1aa73590eee3e78add44a0")
	secondAccount:=common.HexToAddress("58185e1874446969fff65372b404af9a51842156")
	voterStruct,_:=ballot.Voters(nil,addr)
	fmt.Printf("voterStruct:%+v\n", voterStruct)
	printVoters(ballot,firstAccount)
	printVoters(ballot,secondAccount)
	printProposals(ballot)
	ballot.GiveRightToVote(account0Transactor,secondAccount)
	time.Sleep(2000 * time.Millisecond)
	ballot.Vote(account0Transactor,big.NewInt(int64(0)))
	time.Sleep(2000 * time.Millisecond)
	ballot.Vote(account1Transactor,big.NewInt(int64(1)))
	// ballot.Vote(account0Transactor,big.NewInt(int64(2)))
	// ballot.Vote(account0Transactor,big.NewInt(int64(3)))
	time.Sleep(2000 * time.Millisecond)
	printVoters(ballot,firstAccount)
	printVoters(ballot,secondAccount)
	printProposals(ballot)

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
		
	fmt.Println("///////////////////////////////////////")
}