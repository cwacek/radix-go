package radix

import log "github.com/cihub/seelog"
import "sort"

type RadixTreeEntry interface {
  RadixKey() []byte
}

type RadixTree interface {
  Insert(RadixTreeEntry) bool
  Find([]byte) (RadixTreeEntry, bool)
  Len() int
  Walk() []RadixTreeEntry
}

func NewTrie() RadixTree {
  node := new(Trie)
  node.Init()
  return node
}

type ByteSlice []byte

func (this ByteSlice) Len() int { return len(this) }
func (this ByteSlice) Less(i, j int) bool { return this[i] < this[j] }
func (this ByteSlice) Swap(i, j int) {
  tmp := this[i]
  this[i] = this[j]
  this[j] = tmp
}


type Trie struct {
  elemcount int
  root *node
}

type node struct {
  Key byte
  Value RadixTreeEntry
  subtrees map[byte]*node
  elemcount int
}

// Walk a tree, and push all of the elements
// in it into the channel in sorted order
func (n *Trie) Walk() (elems []RadixTreeEntry) {
  elems = make([]RadixTreeEntry, 0, n.elemcount)
  n.root.walk(&elems)
  return elems
}

func (n *node) walk(elemList *[]RadixTreeEntry) {

  if n.Value != nil {
    log.Debugf("Sending value %v", n.Value)
    *elemList = append(*elemList, n.Value)
  }

  if len(n.subtrees) == 0 {
    log.Debugf("Stopping recursion. No sub elements.")
    return
  }

  // Recurse in order
  var keys ByteSlice

  for key := range n.subtrees {
    keys = append(keys, key)
  }

  sort.Sort(keys)

  for _, key := range keys {
    log.Debugf("Recursing to %c", key)
    n.subtrees[key].walk(elemList)
  }
}


func (n *node) Init(key byte) {
  n.subtrees = make(map[byte]*node)
  n.Key = key
  n.Value = nil
}

func (T *Trie) Init() {
  T.elemcount = 0
  T.root = new(node)
  T.root.Init(0)
}

func (T *Trie) Insert(r RadixTreeEntry) (added bool) {

  log.Debugf("Inserting with '%v' at key %s", r, r.RadixKey())

  elem, found := T.root.find(r.RadixKey(), true)

  log.Debugf("Inserting into %v, which has value %v", elem, elem.Value)

  if ! found {
    panic("Couldn't find key when extending")
  }

  if elem.Value == nil {
    /*log.Debugf("'%v' has value. Replacing it with '%v'",elem, r)*/
    elem.Value = r
    added = true
    T.elemcount += 1

  } else {
    elem.Value = r
    added = false
  }
  log.Debugf("elemcount is now %d", T.elemcount)
  return
}

/* Find the element with key 'key'. If extend is true, 
append elements to the tree until the correct one exists. Otherwise return nil
if the element doesn't exist.*/
func (n *node) find(key []byte, extend bool) (elem *node, ok bool) {

  var(
    k byte
    leftover []byte
  )

  k = key[0]
  leftover = key[1:]

  log.Debugf("Subkey is '%s' with len %d", leftover, len(leftover))
  if len(leftover) == 0 && n.Key != 0{
    // This is only the stopping point if we're not at the root.
    // At the root we need to go down one mroe
    log.Debugf("Returning %v", n)
    return n, true
  }

  elem, ok = n.subtrees[k]
  if !ok {
    if ! extend {
      log.Debugf("Didn't find subkey %s", k)
      return nil, false

    } else {
      elem = new(node)
      elem.Init(k)
      log.Debugf("Creating new subnode %v with Key %s", elem, elem.Key)
      n.subtrees[k] = elem
      ok = true
    }
  }

  if len(leftover) > 0 {
    log.Debugf("Searching subtree for key '%s'.", leftover)
    elem, ok = elem.find(leftover, extend)
  }
  return
}

func (n *Trie) Find(key []byte) (RadixTreeEntry, bool) {

  if n.root == nil {
    return nil, false
  }

  log.Debugf("Asked to Find '%s'", key)

  elem, found := n.root.find(key, false)
  log.Debugf("Found %v. Success: %s", elem, found)

  if ! found {
    log.Debugf("Failed")
    return nil, false
  }

  log.Debugf("Found elem %v", elem)
  return elem.Value, true
}

func (T *Trie) Len() int {
  return T.elemcount
}

