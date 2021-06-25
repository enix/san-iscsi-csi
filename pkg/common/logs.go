/*
 * Copyright (c) 2021 Enix, SAS
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 *
 * Authors:
 * Paul Laffitte <paul.laffitte@enix.fr>
 * Alexandre Buisine <alexandre.buisine@enix.fr>
 */

package common

import (
	"context"

	klog "k8s.io/klog/v2"
)

func GetLogKeyAndValues(ctx context.Context, keyAndValues ...interface{}) []interface{} {
	logTags := ctx.Value("logTags").(map[string]interface{})

	for k, v := range logTags {
		keyAndValues = append(keyAndValues, k, v)
	}

	return keyAndValues
}

func LogInfoS(ctx context.Context, msg string, keyAndValues ...interface{}) {
	klog.InfoSDepth(1, msg, GetLogKeyAndValues(ctx, keyAndValues...)...)
}

func AddLogTag(ctx context.Context, key string, value interface{}) {
	logTags := ctx.Value("logTags").(map[string]interface{})
	logTags[key] = value
}
