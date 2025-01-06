package impression

/*
func TestImpressionWithoutData(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/movie/impression", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/movie/impression")

	// テストデータの初期化
	testData := InitImpressionNoData()

	// repository の初期化
	inMemoryRepo := NewInMemoryRepository(testData)
	// サービスの初期化
	queryService := impression.NewQueryService(*inMemoryRepo)
	// ハンドラの初期化
	handler := impression.NewHandler(queryService)

	err := handler.ReaderHandler.GetImpressions(c)
	if err != nil {
		t.Fatal(err)
	}

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedJSON := `[]`

		var expected []*model.MovieImpression
		if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
			t.Fatalf("expected unmarshal is failed")
		}
		var actual []*model.MovieImpression
		if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
			t.Fatalf("actual unmarshal is failed")
		}
		assert.ElementsMatch(t, expected, actual)
	}
}

func TestImpression(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/movie/impression", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/movie/impression")

	// テストデータの初期化
	testData := InitImpression()

	// repository の初期化
	inMemoryRepo := NewInMemoryRepository(testData)
	// サービスの初期化
	queryService := impression.NewQueryService(*inMemoryRepo)
	// ハンドラの初期化
	handler := impression.NewHandler(queryService)

	err := handler.ReaderHandler.GetImpressions(c)
	if err != nil {
		t.Fatal(err)
	}

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedJSON := `[
			{
				"id":1,
				"movie_id":1,
				"status":true,
				"rating":4.5,
				"note":"テスト感想_1",
				"movie":{
					"id":1,
					"title":"テストタイトル_1",
					"overview":"テスト概要_1",
					"release_date":"2024-01-02T03:04:05.000006789+09:00",
					"run_time":123,
					"poster_url":"/poster.jpg",
					"series_id":1,
					"tmdb_id":456,
					"genres":null,
					"poster":null,
					"series":null,
					"movie_impression":null
				},
				"watch_records":[
					{
						"id":1,
						"movie_impression_id":1,
						"watch_media_id":99,
						"watch_date":"2016-12-25T00:00:00+09:00",
						"watch_media":{
							"id":0,
							"code":"",
							"name":null
						}
					},
					{
						"id":2,
						"movie_impression_id":1,
						"watch_media_id":1,
						"watch_date":"2020-01-01T00:00:00+09:00",
						"watch_media":{
							"id":0,
							"code":"",
							"name":null
						}
					}
				]
			},
			{
				"id":2,
				"movie_id":2,
				"status":false,
				"rating":null,
				"note":null,
				"movie":{
					"id":2,
					"title":"テストタイトル_2",
					"overview":null,
					"release_date":"2020-02-03T04:05:06.000000789+09:00",
					"run_time":456,
					"poster_url":null,
					"series_id":null,
					"tmdb_id":null,
					"genres":null,
					"poster":null,
					"series":null,
					"movie_impression":null
				},
				"watch_records":[]
			}
		]`

		var expected []*model.MovieImpression
		if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
			t.Fatalf("expected unmarshal is failed")
		}
		var actual []*model.MovieImpression
		if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
			t.Fatalf("actual unmarshal is failed")
		}
		assert.ElementsMatch(t, expected, actual)
	}
}
*/
