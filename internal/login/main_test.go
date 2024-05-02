package login

// // internal use for tests
// func loadTestDatasetStoredSecrets() []file.DatabaseRecord {
// 	content, err := os.ReadFile("test_data_valid.yml")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	secrets := []file.DatabaseRecord{}
// 	err = yaml.Unmarshal(content, &secrets)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return secrets
// }

// // internal use for tests
// func loadTestDatasetUserInputs() []file.LoginRequest {
// 	content, err := os.ReadFile("test_data_valid.yml")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	user_inputs := []file.LoginRequest{}
// 	err = yaml.Unmarshal(content, &user_inputs)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return user_inputs
// }

// // internal use for tests
// func average_time(series []time.Duration) float64 {
// 	total_time := time.Duration(0.0)
// 	for i := 0; i < len(series); i++ {
// 		total_time += series[i]
// 	}
// 	return float64(total_time) / float64(len(series))
// }

// // internal use for tests
// func standart_deviation_time(series []time.Duration, average float64) float64 {
// 	total_distance := 0.0
// 	for i := 0; i < len(series); i++ {
// 		total_distance += math.Pow(float64(series[i])-average, 2)
// 	}
// 	return math.Sqrt(total_distance / float64(len(series)-1))
// }

// // internal use for tests
// func compare_datasets_on_final_method(stored_secrets []file.DatabaseRecord, invalid_inputs []file.LoginRequest, valid_inputs []file.LoginRequest, test *testing.T) {
// 	if len(valid_inputs) != len(invalid_inputs) || len(valid_inputs) != len(stored_secrets) {
// 		test.Error("Test datasets are incompatible with their lengths.")
// 	}

// 	// Don't count the time spent on cold start
// 	_ = Login(stored_secrets[0], invalid_inputs[0])

// 	// n iterations for valid input
// 	series_valid := []time.Duration{0.0}
// 	for i := 0; i < len(valid_inputs); i++ {

// 		start := time.Now()
// 		_ = Login(stored_secrets[i], valid_inputs[i])
// 		series_valid = append(series_valid, time.Since(start))
// 	}
// 	avg_valid := average_time(series_valid)
// 	std_valid := standart_deviation_time(series_valid, avg_valid)

// 	// n iterations for invalid input
// 	series_invalid := []time.Duration{0.0}
// 	for i := 0; i < len(invalid_inputs); i++ {
// 		start := time.Now()
// 		_ = Login(stored_secrets[i], invalid_inputs[i])
// 		series_invalid = append(series_invalid, time.Since(start))
// 	}
// 	avg_invalid := average_time(series_invalid)
// 	std_invalid := standart_deviation_time(series_invalid, avg_invalid)

// 	avg_ratio := math.Abs(1.0 - avg_invalid/avg_valid)
// 	if avg_ratio > 0.02 {
// 		test.Errorf("Difference between valid and invalid data (%.4f %%) is unacceptable.", avg_ratio*100)
// 	} else {
// 		test.Logf("Difference between valid and invalid data (%.4f %%) is acceptable.", avg_ratio*100)
// 	}

// 	avg_diff := math.Abs(avg_valid-avg_invalid) / 1000000.0 // in milliseconds now
// 	if avg_diff > 2.0 {
// 		test.Errorf("Difference between valid and invalid data (%.4f ms) is unacceptable.", avg_diff)
// 	} else {
// 		test.Logf("Difference between valid and invalid data (%.4f ms) is acceptable.", avg_diff)
// 	}

// 	std_ratio := math.Abs(1.0 - std_invalid/std_valid)
// 	if std_ratio > 0.02 {
// 		test.Errorf("Difference between valid and invalid data (%.4f %%) is unacceptable.", std_ratio*100)
// 	} else {
// 		test.Logf("Difference between valid and invalid data (%.4f %%) is acceptable.", std_ratio*100)
// 	}
// }

// func Test_Username_ValidInput(test *testing.T) {
// 	if !username("ufukty", "ufukty") {
// 		test.Error("Username function didn't accept valid input.")
// 	}
// }
// func Test_Username_InvalidInput(test *testing.T) {
// 	if username("ufukty", "12345") {
// 		test.Error("Username function accepted invalid input.")
// 	}
// }

// func Test_Password_ValidInput(test *testing.T) {
// 	if !Password("$argon2id$v=19$m=32768,t=4,p=2$YjVlYTg3ZGZmNGY0OTdhY2YwMjQ1NjkwNTYxMmNkY2E$iJrKugp+XyB11NA+B3vxwbxlWZiAyDCKkryjLpVtqOQ", "Hello World!") {
// 		test.Error("Password function didn't accept valid input.")
// 	}
// }

// func Test_Password_InvalidInput(test *testing.T) {
// 	if Password("$argon2id$v=19$m=32768,t=4,p=2$YjVlYTg3ZGZmNGY0OTdhY2YwMjQ1NjkwNTYxMmNkY2E$iJrKugp+XyB11NA+B3vxwbxlWZiAyDCKkryjLpVtqOQ", "Hello Moon!") {
// 		test.Error("Password function accepted invalid input.")
// 	}
// }

// func Test_OTP_InvalidInput(test *testing.T) {
// 	if OTP("4JOSKD2SXC4XSUL7DNCTCLX7PRHNDUZW", "111111") {
// 		test.Error("OTP function accepted invalid input.")
// 	}
// }

// func Test_Final_TimingAttack_InvalidUsername(test *testing.T) {
// 	rand.Seed(time.Now().UnixNano())
// 	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 	stored_secret_dataset := loadTestDatasetStoredSecrets()
// 	user_inputs_valid := loadTestDatasetUserInputs()
// 	user_inputs_invalid := loadTestDatasetUserInputs()
// 	for i := 0; i < len(user_inputs_invalid); i++ {
// 		b := make([]rune, 10)
// 		for i := range b {
// 			b[i] = letterRunes[rand.Intn(len(letterRunes))]
// 		}
// 		user_inputs_invalid[i].Username = string(b)
// 	}
// 	compare_datasets_on_final_method(stored_secret_dataset, user_inputs_invalid, user_inputs_valid, test)
// }

// func Test_Final_TimingAttack_EmptyPassword(test *testing.T) {
// 	stored_secret_dataset := loadTestDatasetStoredSecrets()
// 	user_inputs_valid := loadTestDatasetUserInputs()
// 	user_inputs_invalid := loadTestDatasetUserInputs()
// 	for i := 0; i < len(user_inputs_invalid); i++ {
// 		user_inputs_invalid[i].Password = ""
// 	}
// 	compare_datasets_on_final_method(stored_secret_dataset, user_inputs_invalid, user_inputs_valid, test)
// }

// func Test_Final_TimingAttack_InvalidPassword(test *testing.T) {
// 	stored_secret_dataset := loadTestDatasetStoredSecrets()
// 	user_inputs_valid := loadTestDatasetUserInputs()
// 	user_inputs_invalid := loadTestDatasetUserInputs()
// 	for i := 0; i < len(user_inputs_invalid); i++ {
// 		// shuffle every passwords with each other, so, every item in slice will be invalid
// 		user_inputs_invalid[i].Password = user_inputs_invalid[rand.Intn(len(user_inputs_invalid))].Password
// 	}
// 	compare_datasets_on_final_method(stored_secret_dataset, user_inputs_invalid, user_inputs_valid, test)
// }

// func Test_Final_TimingAttack_EmptyOTP(t *testing.T) {
// 	t.Fatal("Test is not implemented yet.")
// }

// func Test_Final_TimingAttack_InvalidOTP(t *testing.T) {
// 	t.Fatal("Test is not implemented yet.")
// }
