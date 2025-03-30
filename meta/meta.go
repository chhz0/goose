package meta

import (
	"time"
)

type ObjectMetaAccessor interface {
	GetObjectMeta() Object
}

type Object interface {
	GetID() uint64
	SetID(id uint64)
	GetName() string
	SetName(name string)
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
	GetUpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time)
}

type ListInterface interface {
	GetTotalCount() int64
	SetTotalCount(total int64)
}

// Type 用以暴露API版本和对象类型
type Type interface {
	GetApiVersion() string
	SetApiVersion(version string)
	GetKind() string
	SetKind(kind string)
}

var _ ListInterface = (*ListMeta)(nil)

// GetTotalCount implements ListInterface.
func (l *ListMeta) GetTotalCount() int64 { return l.TotalCount }

// SetTotalCount implements ListInterface.
func (l *ListMeta) SetTotalCount(total int64) { l.TotalCount = total }

var _ Type = (*TypeMeta)(nil)

// GetApiVersion implements Type.
func (t *TypeMeta) GetApiVersion() string { return t.APIVersion }

// GetKind implements Type.
func (t *TypeMeta) GetKind() string { return t.Kind }

// SetApiVersion implements Type.
func (t *TypeMeta) SetApiVersion(version string) { t.APIVersion = version }

// SetKind implements Type.
func (t *TypeMeta) SetKind(kind string) { t.Kind = kind }

func (o *ListMeta) GetListMeta() ListInterface { return o }

func (o *ObjectMeta) GetObjectMeta() Object { return o }

var _ Object = (*ObjectMeta)(nil)

// GetCreatedAt implements Object.
func (o *ObjectMeta) GetCreatedAt() time.Time { return o.CreatedAt }

// GetID implements Object.
func (o *ObjectMeta) GetID() uint64 { return o.ID }

// GetName implements Object.
func (o *ObjectMeta) GetName() string { return o.Name }

// GetUpdatedAt implements Object.
func (o *ObjectMeta) GetUpdatedAt() time.Time { return o.UpdatedAt }

// SetCreatedAt implements Object.
func (o *ObjectMeta) SetCreatedAt(createdAt time.Time) { o.CreatedAt = createdAt }

// SetID implements Object.
func (o *ObjectMeta) SetID(id uint64) { o.ID = id }

// SetName implements Object.
func (o *ObjectMeta) SetName(name string) { o.Name = name }

// SetUpdatedAt implements Object.
func (o *ObjectMeta) SetUpdatedAt(updatedAt time.Time) { o.UpdatedAt = updatedAt }
