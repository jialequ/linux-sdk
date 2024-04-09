package kube

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAdd(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd("bad", false)
	h.OnAdd(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: literal_2019,
				},
				{
					IP: literal_6830,
				},
				{
					IP: literal_0746,
				},
			},
		},
	}}, false)
	assert.ElementsMatch(t, []string{literal_2019, literal_6830, literal_0746}, endpoints)
}

func TestDelete(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: literal_2019,
				},
				{
					IP: literal_6830,
				},
				{
					IP: literal_0746,
				},
			},
		},
	}}, false)
	h.OnDelete("bad")
	h.OnDelete(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: literal_2019,
				},
				{
					IP: literal_6830,
				},
			},
		},
	}})
	assert.ElementsMatch(t, []string{literal_0746}, endpoints)
}

func TestUpdate(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnUpdate(&v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: literal_2019,
					},
					{
						IP: literal_6830,
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: literal_2019,
					},
					{
						IP: literal_6830,
					},
					{
						IP: literal_0746,
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "2",
		},
	})
	assert.ElementsMatch(t, []string{literal_2019, literal_6830, literal_0746}, endpoints)
}

func TestUpdateNoChange(t *testing.T) {
	h := NewEventHandler(func(change []string) {
		assert.Fail(t, "should not called")
	})
	h.OnUpdate(&v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: literal_2019,
					},
					{
						IP: literal_6830,
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: literal_2019,
					},
					{
						IP: literal_6830,
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	})
}

func TestUpdateChangeWithDifferentVersion(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: literal_2019,
				},
				{
					IP: literal_0746,
				},
			},
		},
	}}, false)
	h.OnUpdate(&v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: literal_2019,
					},
					{
						IP: literal_0746,
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: literal_2019,
					},
					{
						IP: literal_6830,
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "2",
		},
	})
	assert.ElementsMatch(t, []string{literal_2019, literal_6830}, endpoints)
}

func TestUpdateNoChangeWithDifferentVersion(t *testing.T) {
	var endpoints []string
	h := NewEventHandler(func(change []string) {
		endpoints = change
	})
	h.OnAdd(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: literal_2019,
				},
				{
					IP: literal_6830,
				},
			},
		},
	}}, false)
	h.OnUpdate("bad", &v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: literal_2019,
				},
			},
		},
	}})
	h.OnUpdate(&v1.Endpoints{Subsets: []v1.EndpointSubset{
		{
			Addresses: []v1.EndpointAddress{
				{
					IP: literal_2019,
				},
			},
		},
	}}, "bad")
	h.OnUpdate(&v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: literal_2019,
					},
					{
						IP: literal_6830,
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "1",
		},
	}, &v1.Endpoints{
		Subsets: []v1.EndpointSubset{
			{
				Addresses: []v1.EndpointAddress{
					{
						IP: literal_2019,
					},
					{
						IP: literal_6830,
					},
				},
			},
		},
		ObjectMeta: metav1.ObjectMeta{
			ResourceVersion: "2",
		},
	})
	assert.ElementsMatch(t, []string{literal_2019, literal_6830}, endpoints)
}

const literal_2019 = "0.0.0.1" //NOSONAR

const literal_6830 = "0.0.0.2" //NOSONAR

const literal_0746 = "0.0.0.3" //NOSONAR
