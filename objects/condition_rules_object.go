package objects

type FormFieldConditionRules struct {
	ID                     int    `json:"id"`
	FormFieldID            int    `json:"form_field_id"`
	ConditionRuleID        int    `json:"condition_rule_id"`
	Value1                 string `json:"value_1"`
	Value2                 string `json:"value_2"`
	ErrMsg                 string `json:"err_msg"`
	ConditionParentFieldID int    `json:"condition_parent_field_id"`
	ConditionAllRight      bool   `json:"condition_all_right"`
}
