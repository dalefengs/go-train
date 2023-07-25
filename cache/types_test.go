package cache

import (
	"fmt"
	"reflect"
	"testing"
)

type PlayerElement struct {
	Item0 uint32 `tdr_field:"item0" json:"item0"` // 道具1 ID
	Item1 uint32 `tdr_field:"item1" json:"item1"` // 道具2
	Item2 uint32 `tdr_field:"item2" json:"item2"` // 道具3
	Item3 uint32 `tdr_field:"item3" json:"item3"` // 道具4
	Item4 uint32 `tdr_field:"item4" json:"item4"` // 道具5
	Item5 uint32 `tdr_field:"item5" json:"item5"` // 道具6
	Item6 uint32 `tdr_field:"item6" json:"item6"` // 道具7
}

var TermIdMappingFor2v2 = map[uint32]uint32{
	uint32(26486): uint32(223094),
	uint32(26503): uint32(223111),
	uint32(26896): uint32(223504),
	uint32(26900): uint32(223508),
	uint32(28025): uint32(224633),
	uint32(30054): uint32(226662),
	uint32(30085): uint32(226693),
	uint32(26463): uint32(223071),
	uint32(26507): uint32(223115),
	uint32(29427): uint32(226035),
	uint32(30001): uint32(226609),
	uint32(30048): uint32(226656),
	uint32(26481): uint32(223089),
	uint32(26425): uint32(223033),
	uint32(26499): uint32(223107),
	uint32(26614): uint32(223222),
	uint32(30064): uint32(226672),
	uint32(30084): uint32(226692),
	uint32(26569): uint32(223177),
	uint32(26553): uint32(223161),
	uint32(26585): uint32(223193),
	uint32(27206): uint32(223814),
	uint32(28037): uint32(224645),
	uint32(26479): uint32(223087),
	uint32(26428): uint32(223036),
	uint32(26457): uint32(223065),
	uint32(26476): uint32(223084),
	uint32(30083): uint32(226691),
	uint32(30416): uint32(227024),
	uint32(30417): uint32(227025),
	uint32(26412): uint32(223020),
	uint32(26466): uint32(223074),
	uint32(26487): uint32(223095),
	uint32(27140): uint32(223748),
	uint32(30009): uint32(226617),
	uint32(30022): uint32(226630),
	uint32(26442): uint32(223050),
	uint32(28028): uint32(224636),
	uint32(30008): uint32(226616),
	uint32(30047): uint32(226655),
	uint32(30065): uint32(226673),
	uint32(26467): uint32(223075),
	uint32(30087): uint32(226695),
	uint32(30404): uint32(227012),
	uint32(30423): uint32(227031),
	uint32(26434): uint32(223042),
	uint32(30410): uint32(227018),
	uint32(30067): uint32(226675),
	uint32(30418): uint32(227026),
	uint32(30422): uint32(227030),
	uint32(27134): uint32(223742),
	uint32(30057): uint32(226665),
	uint32(30407): uint32(227015),
	uint32(30421): uint32(227029),
	uint32(26527): uint32(223135),
	uint32(26557): uint32(223165),
	uint32(30398): uint32(227006),
	uint32(26534): uint32(223142),
	uint32(26511): uint32(223119),
	uint32(26573): uint32(223181),
	uint32(28021): uint32(224629),
	uint32(26396): uint32(223004),
	uint32(26508): uint32(223116),
	uint32(27397): uint32(224005),
	uint32(30086): uint32(226694),
	uint32(30393): uint32(227001),
	uint32(30405): uint32(227013),
	uint32(26418): uint32(223026),
	uint32(26401): uint32(223009),
	uint32(26470): uint32(223078),
	uint32(26545): uint32(223153),
	uint32(26549): uint32(223157),
	uint32(30406): uint32(227014),
	uint32(30412): uint32(227020),
	uint32(31412): uint32(228020),
	uint32(26398): uint32(223006),
	uint32(26438): uint32(223046),
	uint32(27793): uint32(224401),
	uint32(28029): uint32(224637),
	uint32(28036): uint32(224644),
	uint32(30049): uint32(226657),
	uint32(30056): uint32(226664),
	uint32(30088): uint32(226696),
	uint32(25457): uint32(222065),
	uint32(30420): uint32(227028),
	uint32(26513): uint32(223121),
	uint32(26576): uint32(223184),
	uint32(30063): uint32(226671),
	uint32(30068): uint32(226676),
	uint32(30411): uint32(227019),
	uint32(30415): uint32(227023),
	uint32(26504): uint32(223112),
	uint32(26464): uint32(223072),
	uint32(26494): uint32(223102),
	uint32(30023): uint32(226631),
	uint32(30403): uint32(227011),
	uint32(30409): uint32(227017),
	uint32(26423): uint32(223031),
	uint32(26501): uint32(223109),
	uint32(26516): uint32(223124),
	uint32(26535): uint32(223143),
	uint32(26550): uint32(223158),
	uint32(30424): uint32(227032),
	uint32(26395): uint32(223003),
	uint32(30024): uint32(226632),
	uint32(30425): uint32(227033),
	uint32(26502): uint32(223110),
	uint32(26460): uint32(223068),
	uint32(26477): uint32(223085),
	uint32(26531): uint32(223139),
	uint32(29725): uint32(226333),
	uint32(30401): uint32(227009),
	uint32(30413): uint32(227021),
	uint32(26432): uint32(223040),
	uint32(26483): uint32(223091),
	uint32(26544): uint32(223152),
	uint32(30402): uint32(227010),
	uint32(30419): uint32(227027),
	uint32(26393): uint32(223001),
	uint32(30012): uint32(226620),
	uint32(26445): uint32(223053),
	uint32(26582): uint32(223190),
	uint32(28020): uint32(224628),
	uint32(30394): uint32(227002),
	uint32(26492): uint32(223100),
	uint32(31393): uint32(228001),
	uint32(30408): uint32(227016),
	uint32(30059): uint32(226667),
	uint32(30397): uint32(227005),
	uint32(26403): uint32(223011),
	uint32(26439): uint32(223047),
	uint32(26548): uint32(223156),
	uint32(26577): uint32(223185),
	uint32(30045): uint32(226653),
	uint32(25443): uint32(222051),
}

func Test(t *testing.T) {
	ele := PlayerElement{
		Item0: 26503,
		Item1: 26896,
		Item2: 123123,
		Item3: 123123,
		Item4: 123123,
		Item5: 123123,
		Item6: 26503,
	}
	replaceItemId(&ele)
	fmt.Printf("%v", ele)
}

// 转换截段的 uint32，反射修改字段
func replaceItemIdByReflect(ele *PlayerElement) {
	if ele == nil {
		return
	}
	val := reflect.ValueOf(ele).Elem()
	itemIds := []string{"Item0", "Item1", "Item2", "Item3", "Item4", "Item5", "Item6"}
	for _, itemIds := range itemIds {
		field := val.FieldByName(itemIds)
		if !field.IsValid() { // 无效字段
			continue
		}
		fieldVal := field.Interface().(uint32) // 获取字段值
		// 如果存在映射则修改字段值
		if itemId, ok := TermIdMappingFor2v2[fieldVal]; ok {
			field.Set(reflect.ValueOf(itemId))
		}
	}
}

// 转换截段的 uint32
func replaceItemId(ele *PlayerElement) {
	if ele == nil {
		return
	}
	if itemId, ok := TermIdMappingFor2v2[ele.Item0]; ok {
		ele.Item0 = itemId
	}
	if itemId, ok := TermIdMappingFor2v2[ele.Item1]; ok {
		ele.Item1 = itemId
	}
	if itemId, ok := TermIdMappingFor2v2[ele.Item2]; ok {
		ele.Item2 = itemId
	}
	if itemId, ok := TermIdMappingFor2v2[ele.Item3]; ok {
		ele.Item3 = itemId
	}
	if itemId, ok := TermIdMappingFor2v2[ele.Item4]; ok {
		ele.Item4 = itemId
	}
	if itemId, ok := TermIdMappingFor2v2[ele.Item5]; ok {
		ele.Item5 = itemId
	}
	if itemId, ok := TermIdMappingFor2v2[ele.Item6]; ok {
		ele.Item6 = itemId
	}
}
