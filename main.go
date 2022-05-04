package main

import (
	"fmt"
	"math"
	"sort"
)

const (
	Empty      = -1
	O          = 0
	X          = 1
	InProgress = 2
	Draw       = -2
)

type Node struct {
	State    *State
	Parent   *Node
	Children []*Node
	UCTValue float64
}

func (n *Node) Clone() *Node {
	clonedChildren := make([]*Node, len(n.Children))
	for idx, child := range n.Children {
		clonedChildren[idx] = child.Clone()
	}
	return &Node{
		State:    n.State.Clone(),
		Parent:   n.Parent,
		Children: clonedChildren,
		UCTValue: n.UCTValue,
	}
}

func NewNode(state *State, parent *Node) *Node {
	return &Node{
		State:    state,
		Parent:   parent,
		Children: make([]*Node, 0),
	}
}

func (n *Node) AddChild(child *Node) {
	n.Children = append(n.Children, child)
}

func (n *Node) FillUCTValue() {
	if n.State.VisitCount == 0 {
		n.UCTValue = math.MaxFloat64
		return
	}
	totalVisit := float64(n.Parent.State.VisitCount)
	nodeVisit := float64(n.State.VisitCount)
	n.UCTValue = (n.State.Score / nodeVisit) + 1.41*math.Sqrt(math.Log(totalVisit)/nodeVisit)
}

func (n *Node) CreateChild(move int) {
	nextState := n.State.NextState(move)
	n.AddChild(NewNode(nextState, n))
}

func (n *Node) SelectPromisingNode() *Node {
	if len(n.Children) == 0 {
		return n
	}
	for _, child := range n.Children {
		child.FillUCTValue()
	}
	sort.SliceStable(n.Children, func(i, j int) bool {
		return n.Children[i].UCTValue > n.Children[j].UCTValue
	})
	return n.Children[0]
}

func (n *Node) GetWinnerChild() *Node {
	sort.SliceStable(n.Children, func(i, j int) bool {
		return n.Children[i].State.VisitCount > n.Children[j].State.VisitCount
	})
	return n.Children[0]
}

func (n *Node) ExpandNode() {
	moves := n.State.Board.EmptyPlaces()
	for _, move := range moves {
		n.CreateChild(move)
	}
}

type Tree struct {
	Root *Node
}

func NewTree(board Board, playerNumber int) *Tree {
	rootNode := NewNode(NewState(board, playerNumber), nil)
	// add all possible moves for the given root node
	rootNode.ExpandNode()
	return &Tree{
		Root: rootNode,
	}
}

type Board interface {
	CheckStatus() int
	PerformMove(move, player int)
	EmptyPlaces() []int
	Clone() Board
	Print()
	RandomPlay(player int)
}

type State struct {
	Board        Board
	PlayerNumber int // the last player who did an action
	VisitCount   int
	Score        float64
}

func NewState(board Board, playerNumber int) *State {
	return &State{
		Board:        board,
		PlayerNumber: playerNumber,
	}
}

func (s *State) Clone() *State {
	return &State{
		Board:        s.Board.Clone(),
		PlayerNumber: s.PlayerNumber,
		VisitCount:   s.VisitCount,
		Score:        s.Score,
	}
}

func (s *State) NextState(move int) *State {
	player := s.PlayerNumber ^ 1
	clonedBoard := s.Board.Clone()
	clonedBoard.PerformMove(move, player)
	return NewState(clonedBoard, player)
}

func (s *State) RandomPlay() *State {
	player := s.PlayerNumber ^ 1
	cloneBoard := s.Board.Clone()
	cloneBoard.RandomPlay(player)
	return NewState(cloneBoard, player)
}

type MonteCarloTreeSearch struct {
	Level    int
	Opponent int
}

func (m *MonteCarloTreeSearch) FindNextMove2(board Board, playerNumber int, numIterations int) Board {
	m.Opponent = playerNumber ^ 1
	tree := NewTree(board, m.Opponent)
	iteration := 0
	currentNode := tree.Root
	for iteration < numIterations {
		// is the current node a leaf?
		currentNode = currentNode.SelectPromisingNode()
		if currentNode.State.Board.CheckStatus() != InProgress {
			currentNode = tree.Root
		}
		if len(currentNode.Children) == 0 && currentNode.State.VisitCount == 0 {
			playout := m.SimulateRandomPlayout(currentNode)
			m.BackPropagation(currentNode, playout)
			currentNode = tree.Root
		} else if len(currentNode.Children) == 0 && currentNode.State.VisitCount > 0 {
			currentNode.ExpandNode()
			currentNode = currentNode.Children[0]
			playout := m.SimulateRandomPlayout(currentNode)
			m.BackPropagation(currentNode, playout)
			currentNode = tree.Root
		}
		iteration += 1
	}
	winner := tree.Root.SelectPromisingNode()
	return winner.State.Board
}

func (m *MonteCarloTreeSearch) SimulateRandomPlayout(node *Node) int {
	nextState := node.State.RandomPlay()
	board := nextState.Board
	player := nextState.PlayerNumber ^ 1
	for board.CheckStatus() == InProgress {
		board.RandomPlay(player)
		player = player ^ 1
	}
	return board.CheckStatus()
}

func (m *MonteCarloTreeSearch) BackPropagation(node *Node, playerNumber int) {
	if playerNumber == m.Opponent {
		node.State.Score = math.MinInt
		return
	}
	currNode := node
	for currNode != nil {
		currNode.State.VisitCount += 1
		if currNode.State.PlayerNumber == playerNumber {
			currNode.State.Score += 10
		}
		currNode = currNode.Parent
	}
}

func main() {
	mcts := MonteCarloTreeSearch{
		Opponent: X,
	}
	board := NewTicTacToeBoard()
	player := O
	for board.CheckStatus() == InProgress {
		if player == O {
			var m int
			_, _ = fmt.Scanln(&m)
			board.PerformMove(m, player)
		} else {
			board = mcts.FindNextMove2(board, player, 1000)
		}
		player = player ^ 1
		board.Print()
		fmt.Println("#########")
	}
	fmt.Println("Game over")

}
