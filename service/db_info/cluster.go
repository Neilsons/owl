package db_info

import (
	"time"

	"github.com/ibanyu/owl/service"
	"github.com/ibanyu/owl/util"
)

type OwlCluster struct {
	ID          int64  `json:"id" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Addr        string `json:"addr" gorm:"column:addr"` //ip : port
	User        string `json:"user" gorm:"column:user"`
	Pwd         string `json:"pwd" gorm:"column:pwd"`

	Ct       int64  `json:"ct" gorm:"column:ct"`
	Ut       int64  `json:"ut" gorm:"column:ut"`
	Operator string `json:"operator" gorm:"column:operator"`
}

type ClusterDao interface {
	AddCluster(cluster *OwlCluster) (int64, error)
	UpdateCluster(cluster *OwlCluster) error
	DelCluster(id int64) error
	GetClusterByName(clusterName string) (*OwlCluster, error)
	ListCluster() ([]OwlCluster, error)
}

var clusterDao ClusterDao

func SetClusterDao(impl ClusterDao) {
	clusterDao = impl
}

func AddCluster(cluster *OwlCluster) (int64, error) {
	cryptoData, err := util.AesCrypto([]byte(cluster.Pwd))
	if err != nil {
		return 0, err
	}
	cluster.Ct = time.Now().Unix()

	cluster.Pwd = util.StringifyByteDirectly(cryptoData)
	return clusterDao.AddCluster(cluster)
}

func UpdateCluster(cluster *OwlCluster) error {
	if cluster.Pwd == pwdReplace {
		cluster.Pwd = ""
	}

	if cluster.Pwd != "" {
		cryptoData, err := util.AesCrypto([]byte(cluster.Pwd))
		if err != nil {
			return err
		}

		cluster.Pwd = util.StringifyByteDirectly(cryptoData)
	}
	cluster.Ut = service.Clock.NowUnix()

	return clusterDao.UpdateCluster(cluster)
}

func DelCluster(id int64) error {
	return clusterDao.DelCluster(id)
}

func GetClusterByName(name string) (*OwlCluster, error) {
	cluster, err := clusterDao.GetClusterByName(name)
	if err != nil {
		return nil, err
	}

	deCryptoData, err := util.AesDeCrypto(util.ParseStringedByte(cluster.Pwd))
	if err != nil {
		return nil, err
	}
	cluster.Pwd = string(deCryptoData)
	return cluster, nil
}

func ListCluster() ([]OwlCluster, error) {
	return clusterDao.ListCluster()
}

const pwdReplace = "******"

func ListClusterForUI() ([]OwlCluster, error) {
	clusters, err := ListCluster()
	if err != nil {
		return nil, err
	}

	for i, _ := range clusters {
		clusters[i].Pwd = pwdReplace
	}
	return clusters, nil
}
