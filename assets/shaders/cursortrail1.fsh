#version 330

uniform sampler2DArray tex;
uniform vec4 col_tint;
uniform float points;

in vec2 tex_coord;
in float index;
in vec4 color_pass;

out vec4 color;

void main() {
    vec4 in_color = texture(tex, vec3(tex_coord, 0));
	color = in_color * col_tint * color_pass * vec4(1, 1, 1, 1-smoothstep(points / 3, points, index));
}
