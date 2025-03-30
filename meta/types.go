package meta

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type ExtenAttrs map[string]any

func (e ExtenAttrs) String() string {
	buf, _ := json.Marshal(e)
	return string(buf)
}

func (e ExtenAttrs) Merge(attrsJson string) ExtenAttrs {
	if attrsJson == "" {
		return e
	}

	var newAttrs ExtenAttrs
	_ = json.Unmarshal([]byte(attrsJson), &newAttrs)
	for k, v := range newAttrs {
		if _, ok := e[k]; !ok {
			e[k] = v
		}
	}

	return e
}

// TypeMeta 定义统一的类型元数据结构
type TypeMeta struct {
	// Kind 资源类型
	Kind string `json:"kind,omitempty"`

	// APIVersion 资源版本
	APIVersion string `json:"api_version,omitempty"`
}

type ListMeta struct {
	TotalCount int64 `json:"total_count,omitempty"`
}

// ObjectMeta 定义统一的元数据结构
type ObjectMeta struct {
	// UID 映射到数据库里表的uid字段，不作为资源的唯一标识.
	ID uint64 `json:"id,omitempty" gorm:"primaryKey;column:id;autoIncrement"`

	// InstanceID 实例ID，资源的唯一标识(数据库级别)，一般为"profix-XXXXXX"
	InstanceID string `json:"instanceID,omitempty" gorm:"column:instance_id;type:varchar(32);not null;uniqueIndex:instanceID_UNIQUE"`

	// Name 资源名称，由用户定义
	Name string `json:"name,omitempty" gorm:"column:name;type:varchar(32);not null;uniqueIndex:idx_name_UNIQUE"`

	// Extn 扩展字段，用户自定义
	ExtenAttrs ExtenAttrs `json:"extn,omitempty" gorm:"-"`
	// ExtnShadow 用于存储扩展字段的JSON字符串，避免JSON解析时对字段进行二次解析
	ExtenShadow string `json:"-" gorm:"column:extn_shadow;type:longtext"`

	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at;type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;type:timestamp;autoUpdateTime"`
}

// BeforeCreate 在创建资源时，将扩展字段转换为JSON字符串
func (o *ObjectMeta) BeforeCreate(tx *gorm.DB) (err error) {
	o.ExtenShadow = o.ExtenAttrs.String()

	return nil
}

// BeforeUpdate 在更新资源时，将扩展字段转换为JSON字符串
func (o *ObjectMeta) BeforeUpdate(tx *gorm.DB) (err error) {
	o.ExtenShadow = o.ExtenAttrs.String()

	return nil
}

// AfterFind 在查询资源时，将JSON字符串转换为扩展字段
func (o *ObjectMeta) AfterFind(tx *gorm.DB) (err error) {
	if err := json.Unmarshal([]byte(o.ExtenShadow), &o.ExtenAttrs); err != nil {
		return err
	}

	return nil
}

// ListOptions 定义查询选项
type ListOptions struct {
	TypeMeta `json:",inline"`

	// LabelSelector 用以查找到匹配的 REST 资源.
	LabelSelector string `json:"labelSelector,omitempty"`

	// FieldSelector 用以查找到匹配的 REST 资源, 默认为 所有.
	// 对应数据库中的表字段
	FieldSelector string `json:"fieldSelector,omitempty"`

	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	Offset *int64 `json:"offset,omitempty"`

	Limit *int64 `json:"limit,omitempty"`
}

type ExportOptions struct {
	TypeMeta `json:",inline"`

	Export bool `json:"export"`

	Exact bool `json:"exact"`
}

type GetOptions struct {
	TypeMeta `json:",inline"`
}

type DeleteOptions struct {
	TypeMeta `json:",inline"`

	Unscoped bool `json:"unscoped"`
}

type CreateOptions struct {
	TypeMeta `json:",inline"`

	DryRun []string `json:"dryRun,omitempty"`
}

type PatchOptions struct {
	TypeMeta `json:",inline"`

	DryRun []string `json:"dryRun,omitempty"`

	Force bool `json:"force,omitempty"`
}

type UpdateOptions struct {
	TypeMeta `json:",inline"`

	DryRun []string `json:"dryRun,omitempty"`
}

type AuthorizeOptions struct {
	TypeMeta `json:",inline"`
}

type TableOptions struct {
	TypeMeta `json:",inline"`

	NoHeaders bool `json:"-"`
}
