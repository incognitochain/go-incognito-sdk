package incognitoclient

import "fmt"

func (t *IncognitoTestSuite) TestGetReward() {
	data, err := t.stake.GetRewardAmount("12S4pXkuBjX5sVnhZTUX45DAxv1LFX9oeKicievv58CYNzYDefTHy5Ja3Yiyw2kd1Fx5wQCngX1g7vPe6Q931GgoowDQkgkDqa26jU7")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(data)
}

func (t *IncognitoTestSuite) TestGetNodeValidatorKey() {
	data, err := t.stake.GetStatusNodeValidator("12P3xe6Fnku9NXdXtjqi4rXJG19Cyyx5KbqBSrDsDt2tZrN8oC8")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(data)
}

func (t *IncognitoTestSuite) TestGetTotalStaker() {
	data, err := t.stake.GetTotalStaker()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(data)
}

func (t *IncognitoTestSuite) TestListUnstake() {
	data, _ := t.stake.ListUnstake()
	fmt.Println(data)

	t.NotEmpty(data)
}

func (t *IncognitoTestSuite) TestListRewardAmounts() {
	value, err := t.stake.ListRewardAmounts()

	fmt.Println(value)
	fmt.Println(err)

	t.NotEmpty(value)
}