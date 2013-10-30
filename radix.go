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
  Key []byte
  Value RadixTreeEntry
  subtrees map[byte]*node
  elemcount int
}

// Walk a tree, and push all of the elements
// in it into the channel in sorted order
func (n *Trie) Walk() (elems []RadixTreeEntry) {
  elems = make([]RadixTreeEntry, 0, n.elemcount)
  n.root.walk(&elems, []byte{})
  return elems
}

func (n *node) walk(elemList *[]RadixTreeEntry, currKey []byte) {

  currKey = append(currKey, n.Key...)

  if n.Value != nil {
    log.Debugf("Sending value %#v at %s", n.Value, currKey)
    *elemList = append(*elemList, n.Value)
  }

  if len(n.subtrees) == 0 {
    log.Debugf("Stopping recursion at %s. No sub elements.", currKey)
    return
  }

  // Recurse in order
  var keys ByteSlice

  for key := range n.subtrees {
    keys = append(keys, key)
  }

  sort.Sort(keys)

  for _, key := range keys {
    log.Debugf("Recursing to %s%c", currKey, key)
    n.subtrees[key].walk(elemList, currKey)
  }
}


func (n *node) Init(key []byte) {
  n.subtrees = make(map[byte]*node)
  n.Key = key
  n.Value = nil
}

func (T *Trie) Init() {
  T.elemcount = 0
  T.root = new(node)
  T.root.Init(nil)
  log.Critical("Initializing Trie")
}

func (T *Trie) Insert(r RadixTreeEntry) (added bool) {

  log.Debugf("Inserting with '%#v' at key %s", r, r.RadixKey())

  elem, found := T.root.find(r.RadixKey(), true)

  log.Tracef("Inserting into %v, which has value %v", elem, elem.Value)

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
  log.Tracef("elemcount is now %d", T.elemcount)
  return
}

/* Find the element with key 'key'. If extend is true, 
append elements to the tree until the correct one exists. Otherwise return nil
if the element doesn't exist.*/
func (n *node) find(key []byte, extend bool) (elem *node, ok bool) {

  var(
    k int
    N *node
  )

    N = n
    k = 1
    //Get below the root
    elem, ok = N.subtrees[key[0]]
    switch {
    case ok:
      N = elem

    case !ok && extend:
      elem = new(node)
      elem.Init(key[1:])
      N.subtrees[key[0]] = elem
      N = elem

    case !ok:
      return nil, false
    }

    for ; k < len(key); {
      var i int
      for i = 0; k + i < len(key) && i < len(N.Key) && key[k + i] == N.Key[i]; i++ {
        log.Tracef("%c matches at position %d", key[k+i], i)
      }
      log.Tracef("k=%d, i=%d, key=%s, N.Key=%s", k, i, key[k:], N.Key)

      if i < len(N.Key) {
        //Split this node.
        log.Tracef("Splitting the current node to have key %s with new subtree starting at %v with value %s", N.Key[:i], N.Key[i], N.Key[i:])
        elem = new(node)
        elem.Init(N.Key[i:])
        elem.Value = N.Value
        // Copy the current subtree to the split value
        for subtree, val := range N.subtrees {
          elem.subtrees[subtree] = val
        }

        // Delete the old value *and subkeys*
        N.Value = nil

        N.subtrees = map[byte]*node{N.Key[i]: elem}
        N.Key = N.Key[:i]
      }

      log.Tracef("Subkey is '%s' at %s",
      key[k+i:], N.Key)
      if k + i == len(key) && N.Key != nil {
        // This is only the stopping point if we're not at the root.
        // At the root we need to go down one mroe
        log.Tracef("Returning %#v", N)
        return N, true
      }

      log.Tracef("Looking for %v in subtrees: %v", key[k+i], N.subtrees)

      elem, ok = N.subtrees[key[k+i]]
      if !ok {
        if ! extend {
          log.Tracef("Didn't find subkey %s", k)
          return nil, false

        } else {
          elem = new(node)
          elem.Init(key[k+i:])
          log.Tracef("Creating new subnode %v with Key %s", elem, elem.Key)
          N.subtrees[key[k+i]] = elem
          N = elem
          break
        }
      }

      k += len(N.Key)
      N = elem

    }

    return N, true
}

func (n *Trie) Find(key []byte) (RadixTreeEntry, bool) {

  if n.root == nil {
    return nil, false
  }

  log.Debugf("Asked to Find '%s'", key)

  elem, found := n.root.find(key, false)
  log.Debugf("Found %#v. Success: %s", elem, found)

  if ! found {
    log.Debugf("Failed")
    return nil, false
  }

  log.Debugf("Found elem %#v", elem)
  return elem.Value, true
}

func (T *Trie) Len() int {
  return T.elemcount
}

