package biz

import "github.com/muidea/magicDefault/model"

type Owner2Totalizer map[string]*model.Totalizer

type Type2Totalizer map[int]Owner2Totalizer

type Namespace2Totalizer map[string]Type2Totalizer
