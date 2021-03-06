/*
Copyright 2020 The Nocalhost Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cluster_user

import (
	"context"
	"nocalhost/internal/nocalhost-api/model"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type ClusterUserRepo interface {
	Create(ctx context.Context, model model.ClusterUserModel) (model.ClusterUserModel, error)
	Delete(ctx context.Context, id uint64) error
	DeleteByWhere(ctx context.Context, models model.ClusterUserModel) error
	BatchDelete(ctx context.Context, id []uint64) error
	GetFirst(ctx context.Context, models model.ClusterUserModel) (*model.ClusterUserModel, error)
	GetJoinCluster(ctx context.Context, condition model.ClusterUserJoinCluster) ([]*model.ClusterUserJoinCluster, error)
	GetList(ctx context.Context, models model.ClusterUserModel) ([]*model.ClusterUserModel, error)
	Update(ctx context.Context, models *model.ClusterUserModel) (*model.ClusterUserModel, error)
	UpdateKubeConfig(ctx context.Context, models *model.ClusterUserModel) (*model.ClusterUserModel, error)
	Close()
}

type clusterUserRepo struct {
	db *gorm.DB
}

func NewApplicationClusterRepo(db *gorm.DB) ClusterUserRepo {
	return &clusterUserRepo{
		db: db,
	}
}

func (repo *clusterUserRepo) UpdateKubeConfig(ctx context.Context, models *model.ClusterUserModel) (*model.ClusterUserModel, error) {
	affect := repo.db.Model(&model.ClusterUserModel{}).Where("id=?", models.ID).Update("kubeconfig", models.KubeConfig).RowsAffected
	if affect > 0 {
		return models, nil
	}
	return models, errors.New("update fail")
}

func (repo *clusterUserRepo) GetJoinCluster(ctx context.Context, condition model.ClusterUserJoinCluster) ([]*model.ClusterUserJoinCluster, error) {
	var cluserUserJoinCluster []*model.ClusterUserJoinCluster
	result := repo.db.Table("clusters_users as cluster_user_join_clusters").Select("cluster_user_join_clusters.id,cluster_user_join_clusters.application_id,cluster_user_join_clusters.user_id,cluster_user_join_clusters.cluster_id,cluster_user_join_clusters.namespace,c.name as admin_cluster_name,c.kubeconfig as admin_cluster_kubeconfig").Joins("join clusters as c on cluster_user_join_clusters.cluster_id=c.id").Where(condition).Scan(&cluserUserJoinCluster)
	if result.Error != nil {
		return cluserUserJoinCluster, result.Error
	}
	return cluserUserJoinCluster, nil
}

func (repo *clusterUserRepo) DeleteByWhere(ctx context.Context, models model.ClusterUserModel) error {
	result := repo.db.Unscoped().Delete(models)
	return result.Error
}

func (repo *clusterUserRepo) BatchDelete(ctx context.Context, ids []uint64) error {
	result := repo.db.Unscoped().Delete(model.ClusterUserModel{}, ids)
	if result.RowsAffected > 0 {
		return nil
	}
	return result.Error
}

func (repo *clusterUserRepo) Delete(ctx context.Context, id uint64) error {
	result := repo.db.Unscoped().Delete(model.ClusterUserModel{}, id)
	if result.RowsAffected > 0 {
		return nil
	}
	return result.Error
}

func (repo *clusterUserRepo) Update(ctx context.Context, models *model.ClusterUserModel) (*model.ClusterUserModel, error) {
	where := model.ClusterUserModel{
		ApplicationId: models.ApplicationId,
		ID:            models.ID,
		UserId:        models.UserId,
	}
	_, err := repo.GetFirst(ctx, where)
	if err != nil {
		return models, errors.Wrap(err, "[clsuter_user_repo] get clsuter_user denied")
	}
	emptyModel := model.ClusterUserModel{}
	affectRow := repo.db.Model(&emptyModel).Update(models).RowsAffected
	if affectRow > 0 {
		return models, nil
	}
	return models, errors.Wrap(err, "[clsuter_user_repo] update clsuter_user err")
}

func (repo *clusterUserRepo) GetList(ctx context.Context, models model.ClusterUserModel) ([]*model.ClusterUserModel, error) {
	result := make([]*model.ClusterUserModel, 0)
	repo.db.Where(&models).Find(&result)
	if len(result) > 0 {
		return result, nil
	}
	return nil, errors.New("users cluster not found")
}

func (repo *clusterUserRepo) GetFirst(ctx context.Context, models model.ClusterUserModel) (*model.ClusterUserModel, error) {
	cluster := model.ClusterUserModel{}
	result := repo.db.Where(&models).First(&cluster)
	if result.Error != nil {
		return &cluster, result.Error
	}
	return &cluster, nil
}

func (repo *clusterUserRepo) Create(ctx context.Context, model model.ClusterUserModel) (model.ClusterUserModel, error) {
	err := repo.db.Create(&model).Error
	if err != nil {
		return model, errors.Wrap(err, "[application_cluster_repo] create application_cluster error")
	}

	return model, nil
}

// Close close db
func (repo *clusterUserRepo) Close() {
	repo.db.Close()
}
