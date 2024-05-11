package zeepLife

import "testing"

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2024/5/9 14:31
* @Package:
 */
func Test_ZeppLife_SetSteps(t *testing.T) {
	app := NewZeppLife("1003941268@knownsec.com", "4f4ezha!")
	err := app.SetSteps(7500)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("success set step")
}
