package registry

import (
)

var ModelRegistry = make(map[string]interface{})

func RegisterModel(name string, model interface{}) {
    ModelRegistry[name] = model
}

func GetRegisteredModels() []interface{} {
    models := make([]interface{}, 0, len(ModelRegistry))
    for _, model := range ModelRegistry {
        models = append(models, model)
    }
    return models
}