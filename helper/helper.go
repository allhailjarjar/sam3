package sam

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/eyedeekay/sam3"
	"github.com/eyedeekay/sam3/i2pkeys"
)

func NetListener(name, samaddr, keyspath string) (net.Listener, error) {
	return I2PListener(name, samaddr, keyspath)
}

// I2PListener is a convenience function which takes a SAM tunnel name, a SAM address and a filename.
// If the file contains I2P keys, it will create a service using that address. If the file does not
// exist, keys will be generated and stored in that file.
func I2PListener(name, samaddr, keyspath string) (*sam3.StreamListener, error) {
	log.Printf("Starting and registering I2P service, please wait a couple of minutes...")
	listener, err := I2PStreamSession(name, samaddr, keyspath)

	if keyspath != "" {
		err = ioutil.WriteFile(keyspath+".i2p.public.txt", []byte(listener.Keys().Addr().Base32()), 0644)
		if err != nil {
			log.Fatalf("error storing I2P base32 address in adjacent text file, %s", err)
		}
	}
	return listener.Listen() //, err
}

// I2PStreamSession is a convenience function which returns a sam3.StreamSession instead
// of a sam3.StreamListener. It also takes care of setting a persisitent key on behalf
// of the user.
func I2PStreamSession(name, samaddr, keyspath string) (*sam3.StreamSession, error) {
	log.Printf("Starting and registering I2P session...")
	sam, err := sam3.NewSAM(samaddr)
	if err != nil {
		log.Fatalf("error connecting to SAM to %s: %s", samaddr, err)
	}
	var keys *i2pkeys.I2PKeys
	if keyspath != "" {
		if _, err := os.Stat(keyspath + ".i2p.private"); os.IsNotExist(err) {
			f, err := os.Create(keyspath + ".i2p.private")
			if err != nil {
				log.Fatalf("unable to open I2P keyfile for writing: %s", err)
			}
			defer f.Close()
			tkeys, err := sam.NewKeys()
			if err != nil {
				log.Fatalf("unable to generate I2P Keys, %s", err)
			}
			keys = &tkeys
			err = i2pkeys.StoreKeysIncompat(*keys, f)
			if err != nil {
				log.Fatalf("unable to save newly generated I2P Keys, %s", err)
			}
		} else {
			tkeys, err := i2pkeys.LoadKeys(keyspath + ".i2p.private")
			if err != nil {
				log.Fatalf("unable to load I2P Keys: %e", err)
			}
			keys = &tkeys
		}
	}
	stream, err := sam.NewStreamSession(name, *keys, sam3.Options_Medium)
	return stream, err
}
