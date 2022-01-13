/*
 *    Copyright 2021 chenquan
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package xstream

// Group represents a group collector.
type Group struct {
	s      <-chan interface{}
	f      KeyFunc
	groups map[interface{}][]interface{}
}

// Input implements Collector.
func (g *Group) Input(c <-chan interface{}) {
	g.s = c
}

// GroupBy returns a Group.
func GroupBy(f KeyFunc) *Group {
	return &Group{s: make(chan interface{}), f: f}
}

// Map returns a map.
func (g *Group) Map() map[interface{}][]interface{} {
	if g.groups == nil {
		g.groups = make(map[interface{}][]interface{})
		for item := range g.s {
			key := g.f(item)
			g.groups[key] = append(g.groups[key], item)
		}
	}

	return g.groups
}
