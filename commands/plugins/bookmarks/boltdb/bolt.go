package boltdb

import "github.com/boltdb/bolt"

const bucketName = "bm"

// Bolt is a bookmark storage using the embedded database bolt
type Bolt struct {
	db *bolt.DB
}

// New returns a ready to use Bolt database for storing bookmarks
// the bolt file is located at the path passed
func New(db string) (*Bolt, error) {
	blt, err := bolt.Open(db, 0600, nil)
	if err != nil {
		return nil, err
	}
	bdb := &Bolt{db: blt}
	return bdb, nil
}

// New creates a new bookmark
func (b *Bolt) New(key, value string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return bkt.Put([]byte(key), []byte(value))
	})
	return err
}

// Lookup returns the bookmark at the key, and true if it is present
func (b *Bolt) Lookup(key string) (val string, ok bool) {
	_ = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucketName))
		if bkt == nil {
			return nil // impossible to have a bookmark
		}
		out := bkt.Get([]byte(key))
		if out != nil {
			val = string(out)
			ok = true
		}
		return nil
	})
	return val, ok
}

// Remove deletes the bookmark at key if it exists
func (b *Bolt) Remove(key string) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return bkt.Delete([]byte(key))
	})
	return err
}

// Dump loops through all of the keys and creates a stringmap of
// all of the current bookmarks
func (b *Bolt) Dump() (out map[string]string, err error) {
	out = map[string]string{}
	err = b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucketName))
		if bkt == nil {
			return nil // impossible to have a bookmark
		}
		bkt.ForEach(func(k, v []byte) error {
			out[string(k)] = string(v)
			return nil
		})
		return nil
	})
	return out, err
}
