package bptree


func (self *BpTree) newInternal() (a uint64, err error) {
	return self.new(func(bytes []byte) error {
		_, err := newInternal(bytes, self.meta.keySize)
		return err
	})
}

func (self *BpTree) newLeaf() (a uint64, err error) {
	return self.new(func(bytes []byte) error {
		_, err := newLeaf(bytes, self.meta.keySize)
		return err
	})
}

func (self *BpTree) new(init func([]byte) error) (uint64, error) {
	a, err := self.bf.Allocate()
	if err != nil {
		return 0, err
	}
	err = self.bf.Do(a, 1, func(bytes []byte) error {
		return init(bytes)
	})
	if err != nil {
		return 0, err
	}
	return a, nil
}

func (self *BpTree) doInternal(a uint64, do func(*internal) error) error {
	return self.do(
		a,
		do,
		func(n *leaf) error {
			return Errorf("Unexpected leaf node")
		},
	)
}

func (self *BpTree) doLeaf(a uint64, do func(*leaf) error) error {
	return self.do(
		a,
		func(n *internal) error {
			return Errorf("Unexpected internal node")
		},
		do,
	)
}

func (self *BpTree) do(
	a uint64,
	internalDo func(*internal) error,
	leafDo func(*leaf) error,
) error {
	return self.bf.Do(a, 1, func(bytes []byte) error {
		flags := flag(bytes[0])
		if flags & INTERNAL != 0 {
			n, err := loadInternal(bytes)
			if err != nil {
				return err
			}
			return internalDo(n)
		} else if flags & LEAF != 0 {
			n, err := loadLeaf(bytes)
			if err != nil {
				return err
			}
			return leafDo(n)
		} else {
			return Errorf("Unknown block type")
		}
	})
}

