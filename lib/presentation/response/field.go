package response

import "time"

func FormatTimeField(tf time.Time) string {
	return tf.Format("2006-01-02 15:04:05")
}

func FormatOptionalTimeField(tf *time.Time) *string {
	if tf == nil {
		return nil
	}

	tfStr := FormatTimeField(*tf)

	return &tfStr
}
