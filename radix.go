package radix

import log "github.com/cihub/seelog"

type RadixTree interface {
  Insert([]byte, interface{}) bool
  Find([]byte) (interface{}, bool)
  Len() int
}

func NewRadixTree() RadixTree {
  node := new(radix_tree)
  node.Init()
  return node
}

type radix_tree struct {
  elemcount int
  root *node
}

type node struct {
  Key byte
  Value interface{}
  subtrees map[byte]*node
  elemcount int
}

func (n *node) Init(key byte) {
  n.subtrees = make(map[byte]*node)
  n.Key = key
  n.Value = nil
}

func (T *radix_tree) Init() {
  T.elemcount = 0
  T.root = new(node)
  T.root.Init(0)
}

func (T *radix_tree) Insert(key []byte, val interface{}) (added bool) {

  log.Debugf("Inserting with '%v' at key %s", val, key)

  elem, found := T.root.find(key, true)

  log.Debugf("Inserting into %v, which has value %v", elem, elem.Value)

  if ! found {
    panic("Couldn't find key when extending")
  }

  if elem.Value == nil {
    log.Debugf("'%v' has value. Replacing it with '%v'",elem, val)
    elem.Value = val
    added = true
    T.elemcount += 1

  } else {
    elem.Value = val
    added = false
  }
  log.Debugf("elemcount is now %d", T.elemcount)
  return
}

/* Find the element with key 'key'. If extend is true, 
append elements to the tree until the correct one exists. Otherwise return nil
if the element doesn't exist.*/
func (n *node) find(key []byte, extend bool) (*node, bool) {

  log.Debugf("Called find with key %s at node %v", key, n)
  log.Flush()

  var(
    elem *node
    k byte
    leftover []byte
  )

  k = key[0]
  leftover = key[1:]

  log.Debugf("Subkey is '%s' with len %d", leftover, len(leftover))
  if len(leftover) == 0 {
    log.Debugf("Returning %v", n)
    return n, true
  }

  elem, ok := n.subtrees[k]
  if !ok {
    if ! extend {
      log.Debugf("Didn't find subkey %s in skiplist", k)
      return nil, false

    } else {
      elem = new(node)
      elem.Init(k)
      log.Debugf("Creating new subnode %v with Key %s", elem, elem.Key)
      n.subtrees[k] = elem
    }
  }

  log.Debugf("Searching subtree for key '%s'.", leftover)
  subelem, found := elem.find(leftover, extend)
  log.Debugf("Subtree search found %v. Success: %s", subelem,found)
  return subelem, found
}

func (n *radix_tree) Find(key []byte) (interface {}, bool) {

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

func (T *radix_tree) Len() int {
  return T.elemcount
}

