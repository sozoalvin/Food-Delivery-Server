package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type LlNode struct {
	item *FoodInfo //<- cannot be string
	next *Node
}
type LinkedList struct { //creation of single linked list
	head *LlNode
	size int
}

type Node struct { //for Tries Data structure
	children [39]*Node //prefix tries only have array size of 26 (26 alphabets), but we have an array of 39 sizes. We can effectively store an asrrotement of values from a-z and 0-9 in ANY placement
	isEnd    bool
}

type Trie struct { //for Tries Data Structure
	root *Node
}

type BinaryNode struct {
	item  int
	left  *BinaryNode
	right *BinaryNode
}

type BST struct {
	root *BinaryNode
}

func InitUsernameTrie() *Trie {
	result := &Trie{root: &Node{}}
	return result
}

func InitMyFoodList() *Trie {
	result := &Trie{root: &Node{}}
	return result
}

func (t *Trie) PreInsertTrie(u []string, ch chan string) {
	for _, v := range u {
		t.Insert(v)
	}
	ch <- "Food List Search Auto Complete Database Updated"
}

func (t *Trie) Insert(w string) {

	c := strings.ToLower(w)
	wordLength := len(w)
	r := []rune(c)                    //creates a rune based on c so we can change ASCII value
	currentNode := t.root             //will always start from t.root
	for i := 0; i < wordLength; i++ { //always check each interation according to the length of the word
		switch r[i] {
		case 32:
			r[i] = 123
			break
		case 39:
			r[i] = 124
			break
		case 45:
			r[i] = 125
			break
		case 48: //represents value 0 on the keyboard
			r[i] = 126
			break
		case 49: //represents value 1 on the keyboard
			r[i] = 127
			break
		case 50: //represents value 2 on the keyboard
			r[i] = 128
			break
		case 51: //represents value 3 on the keyboard
			r[i] = 129
			break
		case 52: //represents value 4 on the keyboard
			r[i] = 130
			break
		case 53: //represents value 5 on the keyboard
			r[i] = 131
			break
		case 54: //represents value 6 on the keyboard
			r[i] = 132
			break
		case 55: //represents value 7 on the keyboard
			r[i] = 133
			break
		case 56: //represents value 8 on the keyboard
			r[i] = 134
			break
		case 57: //represents value 9 on the keyboard
			r[i] = 135
			break
		}
		charIndex := r[i] - 'a'
		// fmt.Println("do u see this? ")
		//fmt.Println("me", charIndex)
		if currentNode.children[charIndex] == nil {
			currentNode.children[charIndex] = &Node{}
		}
		currentNode = currentNode.children[charIndex] //currentNode.children is of *Node TYPE. in other words, it's pointining to the next node.
		// assgining this; will SHIFT your current node's position.
	}
	currentNode.isEnd = true
}

func (t *Trie) InsertUser(w string) {

	c := strings.ToLower(w)
	wordLength := len(w)
	r := []rune(c)                    //creates a rune based on c so we can change ASCII value
	currentNode := t.root             //will always start from t.root
	for i := 0; i < wordLength; i++ { //always check each interation according to the length of the word
		switch r[i] {
		case 32:
			r[i] = 123
			break
		case 39:
			r[i] = 124
			break
		case 45:
			r[i] = 125
			break
		case 48: //represents value 0 on the keyboard
			r[i] = 126
			break
		case 49: //represents value 1 on the keyboard
			r[i] = 127
			break
		case 50: //represents value 2 on the keyboard
			r[i] = 128
			break
		case 51: //represents value 3 on the keyboard
			r[i] = 129
			break
		case 52: //represents value 4 on the keyboard
			r[i] = 130
			break
		case 53: //represents value 5 on the keyboard
			r[i] = 131
			break
		case 54: //represents value 6 on the keyboard
			r[i] = 132
			break
		case 55: //represents value 7 on the keyboard
			r[i] = 133
			break
		case 56: //represents value 8 on the keyboard
			r[i] = 134
			break
		case 57: //represents value 9 on the keyboard
			r[i] = 135
			break
		}
		charIndex := r[i] - 'a'
		// fmt.Println("me", charIndex)
		if currentNode.children[charIndex] == nil {
			currentNode.children[charIndex] = &Node{}
		}
		currentNode = currentNode.children[charIndex] //currentNode.children is of *Node TYPE. in other words, it's pointining to the next node.
		// assgining this; will SHIFT your current node's position.
	}
	currentNode.isEnd = true
	// ch <- "Userdata base completed"
}

func (t *Trie) Search(w string) bool {

	c := strings.ToLower(w)
	wordLength := len(c)
	currentNode := t.root
	for i := 0; i < wordLength; i++ {
		charIndex := c[i] - 'a'
		if currentNode.children[charIndex] == nil {
			return false
		}
		currentNode = currentNode.children[charIndex]
	}
	if currentNode.isEnd == true {
		return true
	}
	return false
}
func (t *Trie) UserSearch(w string) bool { //to check password

	c := strings.ToLower(w)
	r := []rune(c)
	wordLength := len(c)
	currentNode := t.root
	for i := 0; i < wordLength; i++ {

		switch r[i] {
		case 32:
			r[i] = 123
			break
		case 39:
			r[i] = 124
			break
		case 45:
			r[i] = 125
			break
		case 48: //represents value 0 on the keyboard
			r[i] = 126
			break
		case 49: //represents value 1 on the keyboard
			r[i] = 127
			break
		case 50: //represents value 2 on the keyboard
			r[i] = 128
			break
		case 51: //represents value 3 on the keyboard
			r[i] = 129
			break
		case 52: //represents value 4 on the keyboard
			r[i] = 130
			break
		case 53: //represents value 5 on the keyboard
			r[i] = 131
			break
		case 54: //represents value 6 on the keyboard
			r[i] = 132
			break
		case 55: //represents value 7 on the keyboard
			r[i] = 133
			break
		case 56: //represents value 8 on the keyboard
			r[i] = 134
			break
		case 57: //represents value 9 on the keyboard
			r[i] = 135
			break
		}
		charIndex := r[i] - 'a'
		if currentNode.children[charIndex] == nil {
			return false
		}
		currentNode = currentNode.children[charIndex]
	}
	if currentNode.isEnd == true {
		return true
	}
	//}
	return false
}

func InitFoodDBTrie() *Trie {
	result := &Trie{root: &Node{}}
	return result
}

func InitPostalCode() *BST { //creates POSTAL CODE DB
	result := &BST{nil}
	return result
}

func (bst *BST) insertNode(t **BinaryNode, v int) error {
	if *t == nil {
		newNode := &BinaryNode{
			item:  v,
			left:  nil,
			right: nil,
		}
		*t = newNode

		return nil
	}
	if v < (*t).item {
		bst.insertNode(&((*t).left), v)
	} else {
		bst.insertNode(&((*t).right), v)
	}
	return nil
}

func (bst *BST) PreInsertPostalCode() {
	for _, v := range V {
		bst.insertNode(&bst.root, v.PostalCode)
	}
}

func (bst *BST) searchNode(t *BinaryNode, v int) *BinaryNode {
	if t == nil {
		return nil
	} else {
		if t.item == v {
			return t
		} else {
			if v < t.item {
				return bst.searchNode(t.left, v)
			} else {
				return bst.searchNode(t.right, v)
			}
		}
	}
}

func (bst *BST) Search(v int) (*BinaryNode, error) {

	result := strconv.Itoa(v)
	if len(result) != 6 {
		return nil, errors.New("Postal Code Must Contain 6 digits")
	}

	return bst.searchNode(bst.root, v), nil
}

func (t *Trie) GetSuggestion(query string, total int) []string { //edit

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("\nRecovery from Panic. Please DO NOT use any other characters apart from small case alphabets a-z and numbers 0-9 thank you", r)
			fmt.Println("\nPlease resume running the program\n")
		}
	}()

	var result []string
	//move to next position node from the searching character
	currentNode := t.root //starts from root.

	r := []rune(query)
	for i := 0; i < len(query); i++ {
		switch r[i] {
		case 32: //represents value space on the keyboard
			r[i] = 123
			break
		case 39:
			r[i] = 124 //represents value ' on the keyboard
			break
		case 45: // represents value - on the keyboard
			r[i] = 125
			break
		case 48: //represents value 0 on the keyboard
			r[i] = 126
			break
		case 49: //represents value 1 on the keyboard
			r[i] = 127
			break
		case 50: //represents value 2 on the keyboard
			r[i] = 128
			break
		case 51: //represents value 3 on the keyboard
			r[i] = 129
			break
		case 52: //represents value 4 on the keyboard
			r[i] = 130
			break
		case 53: //represents value 5 on the keyboard
			r[i] = 131
			break
		case 54: //represents value 6 on the keyboard
			r[i] = 132
			break
		case 55: //represents value 7 on the keyboard
			r[i] = 133
			break
		case 56: //represents value 8 on the keyboard
			r[i] = 134
			break
		case 57: //represents value 9 on the keyboard
			r[i] = 135
			break
		}
		charIndex := r[i] - 'a'
		if currentNode.children[charIndex] == nil {
			return result
		}
		currentNode = currentNode.children[charIndex] //helps to move to the next Node.
	}

	if currentNode.isEnd && isLastNode(currentNode) { //if the current Node is the end

		result = append(result, query)
		return result
	}

	wordList := []string{}

	if currentNode.isEnd {
		wordList = append(wordList, query)
		total--
	}

	if !isLastNode(currentNode) { //this returns a true if the current Node is pointing to other children nodes. So if no more children then correct; just return false
		_, result = Suggestion(query, wordList, total, currentNode)
	}

	return result
}

func Suggestion(prefix string, wordList []string, repeat int, currentNode *Node) (int, []string) { //edit

	if isLastNode(currentNode) { //if there are children nodes

		if currentNode.isEnd && len(wordList) < 1 { //your current word can be an end; and it can also point to other children
			wordList = append(wordList, prefix)
		}

		return repeat, wordList

	}

	for i := 0; i < 39; i++ {
		if repeat < 1 {
			return repeat, wordList
		}
		nt := currentNode
		if currentNode.children[i] != nil { //this checks if there are any children nodes have any values.
			newl := ConvItoS(i) //function that gives int value but comes back with alphabetical representation of our customized tree struct
			prefix += string(newl)
			nt = nt.children[i]

			if nt.isEnd {
				wordList = append(wordList, prefix)
				repeat--
			}

			repeat, wordList = Suggestion(prefix, wordList, repeat, nt)

			prefix = prefix[0 : len(prefix)-1]
		}

	}

	return repeat, wordList

}

func isLastNode(nextNode *Node) bool { //this function checks for his friends.

	for i := 0; i < 39; i++ {
		if nextNode.children[i] != nil {
			return false
		}
	}

	return true

}

func ConvItoS(i int) string {

	array := [39]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", " ", "'", "-", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	return array[i]

}
