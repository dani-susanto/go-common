package test

func GetValidationErrorTest(data any) map[string][]string {
	validationErrors := make(map[string][]string)

	fields, ok := data.(map[string]any)
	if !ok {
		return validationErrors
	}

	for field, value := range fields {
		messages, ok := value.([]any)
		if !ok {
			continue
		}

		for _, message := range messages {
			text, ok := message.(string)
			if !ok {
				continue
			}

			validationErrors[field] = append(
				validationErrors[field],
				text,
			)
		}
	}

	return validationErrors
}
