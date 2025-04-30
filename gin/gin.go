package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/saneechka/serializer"
)

func MyBindJSON(c *gin.Context, obj any) error {
	s, err := serializer.New("json")
	if err != nil {
		return err
	}

	data, err := c.GetRawData()
	if err != nil {
		return err
	}

	return s.Unmarshal(data, obj)
}

func MyBindTOML(c *gin.Context, obj any) error {
	s, err := serializer.New("toml")
	if err != nil {
		return err
	}

	data, err := c.GetRawData()
	if err != nil {
		return err
	}

	return s.Unmarshal(data, obj)
}

func MyJSON(c *gin.Context, code int, obj any) error {
	s, err := serializer.New("json")
	if err != nil {
		return err
	}

	data, err := s.Marshal(obj)
	if err != nil {
		return err
	}

	c.Data(code, binding.MIMEJSON, data)
	return nil
}

func MyTOML(c *gin.Context, code int, obj any) error {
	s, err := serializer.New("toml")
	if err != nil {
		return err
	}

	data, err := s.Marshal(obj)
	if err != nil {
		return err
	}

	c.Data(code, "application/toml", data)
	return nil
}
