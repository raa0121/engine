// File generated by G3NSHADERS. Do not edit.
// To regenerate this file install 'g3nshaders' and execute:
// 'go generate' in this folder.

package shaders

const include_attributes_source = `//
// Vertex attributes
//
layout(location = 0) in  vec3  VertexPosition;
layout(location = 1) in  vec3  VertexNormal;
layout(location = 2) in  vec3  VertexColor;
layout(location = 3) in  vec2  VertexTexcoord;
layout(location = 4) in  float VertexDistance;
layout(location = 5) in  vec4  VertexTexoffsets;


`

const include_lights_source = `//
// Lights uniforms
//

// Ambient lights uniforms
#if AMB_LIGHTS>0
    uniform vec3 AmbientLightColor[AMB_LIGHTS];
#endif

// Directional lights uniform array. Each directional light uses 2 elements
#if DIR_LIGHTS>0
    uniform vec3 DirLight[2*DIR_LIGHTS];
    // Macros to access elements inside the DirectionalLight uniform array
    #define DirLightColor(a)		DirLight[2*a]
    #define DirLightPosition(a)		DirLight[2*a+1]
#endif

// Point lights uniform array. Each point light uses 3 elements
#if POINT_LIGHTS>0
    uniform vec3 PointLight[3*POINT_LIGHTS];
    // Macros to access elements inside the PointLight uniform array
    #define PointLightColor(a)			PointLight[3*a]
    #define PointLightPosition(a)		PointLight[3*a+1]
    #define PointLightLinearDecay(a)	PointLight[3*a+2].x
    #define PointLightQuadraticDecay(a)	PointLight[3*a+2].y
#endif

#if SPOT_LIGHTS>0
    // Spot lights uniforms. Each spot light uses 5 elements
    uniform vec3  SpotLight[5*SPOT_LIGHTS];
    
    // Macros to access elements inside the PointLight uniform array
    #define SpotLightColor(a)			SpotLight[5*a]
    #define SpotLightPosition(a)		SpotLight[5*a+1]
    #define SpotLightDirection(a)		SpotLight[5*a+2]
    #define SpotLightAngularDecay(a)	SpotLight[5*a+3].x
    #define SpotLightCutoffAngle(a)		SpotLight[5*a+3].y
    #define SpotLightLinearDecay(a)		SpotLight[5*a+3].z
    #define SpotLightQuadraticDecay(a)	SpotLight[5*a+4].x
#endif

`

const include_material_source = `//
// Material properties uniform
//

// Material parameters uniform array
uniform vec3 Material[6];
// Macros to access elements inside the Material array
#define MatAmbientColor		Material[0]
#define MatDiffuseColor     Material[1]
#define MatSpecularColor    Material[2]
#define MatEmissiveColor    Material[3]
#define MatShininess        Material[4].x
#define MatOpacity          Material[4].y
#define MatPointSize        Material[4].z
#define MatPointRotationZ   Material[5].x

#if MAT_TEXTURES > 0
    // Texture unit sampler array
    uniform sampler2D MatTexture[MAT_TEXTURES];
    // Texture parameters (3*vec2 per texture)
    uniform vec2 MatTexinfo[3*MAT_TEXTURES];
    // Macros to access elements inside the MatTexinfo array
    #define MatTexOffset(a)		MatTexinfo[(3*a)]
    #define MatTexRepeat(a)		MatTexinfo[(3*a)+1]
    #define MatTexFlipY(a)		bool(MatTexinfo[(3*a)+2].x)
    #define MatTexVisible(a)	bool(MatTexinfo[(3*a)+2].y)
#endif

// GLSL 3.30 does not allow indexing texture sampler with non constant values.
// This macro is used to mix the texture with the specified index with the material color.
// It should be called for each texture index. It uses two externally defined variables:
// vec4 texColor
// vec4 texMixed
#define MIX_TEXTURE(i)                                                                       \
    if (MatTexVisible(i)) {                                                                  \
        texColor = texture(MatTexture[i], FragTexcoord * MatTexRepeat(i) + MatTexOffset(i)); \
        if (i == 0) {                                                                        \
            texMixed = texColor;                                                             \
        } else {                                                                             \
            texMixed = mix(texMixed, texColor, texColor.a);                                  \
        }                                                                                    \
    }

`

const include_phong_model_source = `/***
 phong lighting model
 Parameters:
    position:   input vertex position in camera coordinates
    normal:     input vertex normal in camera coordinates
    camDir:     input camera directions
    matAmbient: input material ambient color
    matDiffuse: input material diffuse color
    ambdiff:    output ambient+diffuse color
    spec:       output specular color
 Uniforms:
    AmbientLightColor[]
    DiffuseLightColor[]
    DiffuseLightPosition[]
    PointLightColor[]
    PointLightPosition[]
    PointLightLinearDecay[]
    PointLightQuadraticDecay[]
    MatSpecularColor
    MatShininess
*****/
void phongModel(vec4 position, vec3 normal, vec3 camDir, vec3 matAmbient, vec3 matDiffuse, out vec3 ambdiff, out vec3 spec) {

    vec3 ambientTotal  = vec3(0.0);
    vec3 diffuseTotal  = vec3(0.0);
    vec3 specularTotal = vec3(0.0);

#if AMB_LIGHTS>0
    // Ambient lights
    for (int i = 0; i < AMB_LIGHTS; i++) {
        ambientTotal += AmbientLightColor[i] * matAmbient;
    }
#endif

#if DIR_LIGHTS>0
    // Directional lights
    for (int i = 0; i < DIR_LIGHTS; i++) {
        // Diffuse reflection
        // DirLightPosition is the direction of the current light
        vec3 lightDirection = normalize(DirLightPosition(i));
        // Calculates the dot product between the light direction and this vertex normal.
        float dotNormal = max(dot(lightDirection, normal), 0.0);
        diffuseTotal += DirLightColor(i) * matDiffuse * dotNormal;
        // Specular reflection
        // Calculates the light reflection vector
        vec3 ref = reflect(-lightDirection, normal);
        if (dotNormal > 0.0) {
            specularTotal += DirLightColor(i) * MatSpecularColor * pow(max(dot(ref, camDir), 0.0), MatShininess);
        }
    }
#endif

#if POINT_LIGHTS>0
    // Point lights
    for (int i = 0; i < POINT_LIGHTS; i++) {
        // Common calculations
        // Calculates the direction and distance from the current vertex to this point light.
        vec3 lightDirection = PointLightPosition(i) - vec3(position);
        float lightDistance = length(lightDirection);
        // Normalizes the lightDirection
        lightDirection = lightDirection / lightDistance;
        // Calculates the attenuation due to the distance of the light
        float attenuation = 1.0 / (1.0 + PointLightLinearDecay(i) * lightDistance +
            PointLightQuadraticDecay(i) * lightDistance * lightDistance);
        // Diffuse reflection
        float dotNormal = max(dot(lightDirection, normal), 0.0);
        diffuseTotal += PointLightColor(i) * matDiffuse * dotNormal * attenuation;
        // Specular reflection
        // Calculates the light reflection vector
        vec3 ref = reflect(-lightDirection, normal);
        if (dotNormal > 0.0) {
            specularTotal += PointLightColor(i) * MatSpecularColor *
                pow(max(dot(ref, camDir), 0.0), MatShininess) * attenuation;
        }
    }
#endif

#if SPOT_LIGHTS>0
    for (int i = 0; i < SPOT_LIGHTS; i++) {
        // Calculates the direction and distance from the current vertex to this spot light.
        vec3 lightDirection = SpotLightPosition(i) - vec3(position);
        float lightDistance = length(lightDirection);
        lightDirection = lightDirection / lightDistance;

        // Calculates the attenuation due to the distance of the light
        float attenuation = 1.0 / (1.0 + SpotLightLinearDecay(i) * lightDistance +
            SpotLightQuadraticDecay(i) * lightDistance * lightDistance);

        // Calculates the angle between the vertex direction and spot direction
        // If this angle is greater than the cutoff the spotlight will not contribute
        // to the final color.
        float angle = acos(dot(-lightDirection, SpotLightDirection(i)));
        float cutoff = radians(clamp(SpotLightCutoffAngle(i), 0.0, 90.0));

        if (angle < cutoff) {
            float spotFactor = pow(dot(-lightDirection, SpotLightDirection(i)), SpotLightAngularDecay(i));

            // Diffuse reflection
            float dotNormal = max(dot(lightDirection, normal), 0.0);
            diffuseTotal += SpotLightColor(i) * matDiffuse * dotNormal * attenuation * spotFactor;

            // Specular reflection
            vec3 ref = reflect(-lightDirection, normal);
            if (dotNormal > 0.0) {
                specularTotal += SpotLightColor(i) * MatSpecularColor * pow(max(dot(ref, camDir), 0.0), MatShininess) * attenuation * spotFactor;
            }
        }
    }
#endif

    // Sets output colors
    ambdiff = ambientTotal + MatEmissiveColor + diffuseTotal;
    spec = specularTotal;
}


`

const basic_fragment_source = `//
// Fragment Shader template
//

in vec3 Color;
out vec4 FragColor;

void main() {

    FragColor = vec4(Color, 1.0);
}

`

const basic_vertex_source = `//
// Vertex shader basic
//
#include <attributes>

// Model uniforms
uniform mat4 MVP;

// Final output color for fragment shader
out vec3 Color;

void main() {

    Color = VertexColor;
    gl_Position = MVP * vec4(VertexPosition, 1.0);
}


`

const panel_fragment_source = `//
// Fragment Shader template
//

// Texture uniforms
uniform sampler2D	MatTexture;
uniform vec2		MatTexinfo[3];

// Macros to access elements inside the MatTexinfo array
#define MatTexOffset		MatTexinfo[0]
#define MatTexRepeat		MatTexinfo[1]
#define MatTexFlipY	    	bool(MatTexinfo[2].x) // not used
#define MatTexVisible	    bool(MatTexinfo[2].y) // not used

// Inputs from vertex shader
in vec2 FragTexcoord;

// Input uniform
uniform vec4 Panel[8];
#define Bounds			Panel[0]		  // panel bounds in texture coordinates
#define Border			Panel[1]		  // panel border in texture coordinates
#define Padding			Panel[2]		  // panel padding in texture coordinates
#define Content			Panel[3]		  // panel content area in texture coordinates
#define BorderColor		Panel[4]		  // panel border color
#define PaddingColor	Panel[5]		  // panel padding color
#define ContentColor	Panel[6]		  // panel content color
#define TextureValid	bool(Panel[7].x)  // texture valid flag

// Output
out vec4 FragColor;


/***
* Checks if current fragment texture coordinate is inside the
* supplied rectangle in texture coordinates:
* rect[0] - position x [0,1]
* rect[1] - position y [0,1]
* rect[2] - width [0,1]
* rect[3] - height [0,1]
*/
bool checkRect(vec4 rect) {

    if (FragTexcoord.x < rect[0]) {
        return false;
    }
    if (FragTexcoord.x > rect[0] + rect[2]) {
        return false;
    }
    if (FragTexcoord.y < rect[1]) {
        return false;
    }
    if (FragTexcoord.y > rect[1] + rect[3]) {
        return false;
    }
    return true;
}


void main() {

    // Discard fragment outside of received bounds
    // Bounds[0] - xmin
    // Bounds[1] - ymin
    // Bounds[2] - xmax
    // Bounds[3] - ymax
    if (FragTexcoord.x <= Bounds[0] || FragTexcoord.x >= Bounds[2]) {
        discard;
    }
    if (FragTexcoord.y <= Bounds[1] || FragTexcoord.y >= Bounds[3]) {
        discard;
    }

    // Check if fragment is inside content area
    if (checkRect(Content)) {

        // If no texture, the color will be the material color.
        vec4 color = ContentColor;

		if (TextureValid) {
            // Adjust texture coordinates to fit texture inside the content area
            vec2 offset = vec2(-Content[0], -Content[1]);
            vec2 factor = vec2(1/Content[2], 1/Content[3]);
            vec2 texcoord = (FragTexcoord + offset) * factor;
            vec4 texColor = texture(MatTexture, texcoord * MatTexRepeat + MatTexOffset);

            // Mix content color with texture color.
            // Note that doing a simple linear interpolation (e.g. using mix()) is not correct!
            // The right formula can be found here: https://en.wikipedia.org/wiki/Alpha_compositing#Alpha_blending
            // For a more in-depth discussion: http://apoorvaj.io/alpha-compositing-opengl-blending-and-premultiplied-alpha.html#toc4

            // Pre-multiply the content color
            vec4 contentPre = ContentColor;
            contentPre.rgb *= contentPre.a;

            // Pre-multiply the texture color
            vec4 texPre = texColor;
            texPre.rgb *= texPre.a;

            // Combine colors the premultiplied final color
            color = texPre + contentPre * (1 - texPre.a);

            // Un-pre-multiply (pre-divide? :P)
            color.rgb /= color.a;
		}

        FragColor = color;
        return;
    }

    // Checks if fragment is inside paddings area
    if (checkRect(Padding)) {
        FragColor = PaddingColor;
        return;
    }

    // Checks if fragment is inside borders area
    if (checkRect(Border)) {
        FragColor = BorderColor;
        return;
    }

    // Fragment is in margins area (always transparent)
    FragColor = vec4(1,1,1,0);
}

`

const panel_vertex_source = `//
// Vertex shader panel
//
#include <attributes>

// Model uniforms
uniform mat4 ModelMatrix;

// Outputs for fragment shader
out vec2 FragTexcoord;


void main() {

    // Always flip texture coordinates
    vec2 texcoord = VertexTexcoord;
    texcoord.y = 1 - texcoord.y;
    FragTexcoord = texcoord;

    // Set position
    vec4 pos = vec4(VertexPosition.xyz, 1);
    gl_Position = ModelMatrix * pos;
}

`

const phong_fragment_source = `//
// Fragment Shader template
//

// Inputs from vertex shader
in vec4 Position;       // Vertex position in camera coordinates.
in vec3 Normal;         // Vertex normal in camera coordinates.
in vec3 CamDir;         // Direction from vertex to camera
in vec2 FragTexcoord;

#include <lights>
#include <material>
#include <phong_model>

// Final fragment color
out vec4 FragColor;

void main() {

    // Mix material color with textures colors
    vec4 texMixed = vec4(1);
    vec4 texColor;
    #if MAT_TEXTURES==1
        MIX_TEXTURE(0)
    #elif MAT_TEXTURES==2
        MIX_TEXTURE(0)
        MIX_TEXTURE(1)
    #elif MAT_TEXTURES==3
        MIX_TEXTURE(0)
        MIX_TEXTURE(1)
        MIX_TEXTURE(2)
    #endif

    // Combine material with texture colors
    vec4 matDiffuse = vec4(MatDiffuseColor, MatOpacity) * texMixed;
    vec4 matAmbient = vec4(MatAmbientColor, MatOpacity) * texMixed;

    // Inverts the fragment normal if not FrontFacing
    vec3 fragNormal = Normal;
    if (!gl_FrontFacing) {
        fragNormal = -fragNormal;
    }

    // Calculates the Ambient+Diffuse and Specular colors for this fragment using the Phong model.
    vec3 Ambdiff, Spec;
    phongModel(Position, fragNormal, CamDir, vec3(matAmbient), vec3(matDiffuse), Ambdiff, Spec);

    // Final fragment color
    FragColor = min(vec4(Ambdiff + Spec, matDiffuse.a), vec4(1.0));
}

`

const phong_vertex_source = `//
// Vertex Shader
//
#include <attributes>

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

#include <material>

// Output variables for Fragment shader
out vec4 Position;
out vec3 Normal;
out vec3 CamDir;
out vec2 FragTexcoord;

void main() {

    // Transform this vertex position to camera coordinates.
    Position = ModelViewMatrix * vec4(VertexPosition, 1.0);

    // Transform this vertex normal to camera coordinates.
    Normal = normalize(NormalMatrix * VertexNormal);

    // Calculate the direction vector from the vertex to the camera
    // The camera is at 0,0,0
    CamDir = normalize(-Position.xyz);

    // Flips texture coordinate Y if requested.
    vec2 texcoord = VertexTexcoord;
#if MAT_TEXTURES>0
    if (MatTexFlipY(0)) {
        texcoord.y = 1 - texcoord.y;
    }
#endif
    FragTexcoord = texcoord;

    gl_Position = MVP * vec4(VertexPosition, 1.0);
}

`

const physical_fragment_source = `//
// Physical material fragment shader
//

// Inputs from vertex shader
in vec4 Position;       // Vertex position in camera coordinates.
in vec3 Normal;         // Vertex normal in camera coordinates.
in vec3 CamDir;         // Direction from vertex to camera
in vec2 FragTexcoord;

// Material parameters uniform array
uniform vec4 Material[3];
// Macros to access elements inside the Material array
#define uBaseColor		    Material[0]
#define uEmissiveColor      Material[1]
#define uMetallicFactor     Material[2].x
#define uRoughnessFactor    Material[2].y

#include <lights>

// Final fragment color
out vec4 FragColor;

void main() {


    // Inverts the fragment normal if not FrontFacing
    vec3 fragNormal = Normal;
    if (!gl_FrontFacing) {
        fragNormal = -fragNormal;
    }


    // Final fragment color
    FragColor = uBaseColor;
}


`

const physical_vertex_source = `//
// Physical maiterial vertex shader
//
#include <attributes>

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

// Output variables for Fragment shader
out vec4 Position;
out vec3 Normal;
out vec3 CamDir;
out vec2 FragTexcoord;

void main() {

    // Transform this vertex position to camera coordinates.
    Position = ModelViewMatrix * vec4(VertexPosition, 1.0);

    // Transform this vertex normal to camera coordinates.
    Normal = normalize(NormalMatrix * VertexNormal);

    // Calculate the direction vector from the vertex to the camera
    // The camera is at 0,0,0
    CamDir = normalize(-Position.xyz);

//    // Flips texture coordinate Y if requested.
   vec2 texcoord = VertexTexcoord;
//#if MAT_TEXTURES>0
//    if (MatTexFlipY(0)) {
//        texcoord.y = 1 - texcoord.y;
//    }
//#endif
    FragTexcoord = texcoord;

    gl_Position = MVP * vec4(VertexPosition, 1.0);
}


`

const point_fragment_source = `#include <material>

// GLSL 3.30 does not allow indexing texture sampler with non constant values.
// This macro is used to mix the texture with the specified index with the material color.
// It should be called for each texture index.
#define MIX_POINT_TEXTURE(i)                                                                                     \
    if (MatTexVisible(i)) {                                                                                      \
        vec2 pt = gl_PointCoord - vec2(0.5);                                                                     \
        vec4 texColor = texture(MatTexture[i], (Rotation * pt + vec2(0.5)) * MatTexRepeat(i) + MatTexOffset(i)); \
        if (i == 0) {                                                                                            \
            texMixed = texColor;                                                                                 \
        } else {                                                                                                 \
            texMixed = mix(texMixed, texColor, texColor.a);                                                      \
        }                                                                                                        \
    }

// Inputs from vertex shader
in vec3 Color;
flat in mat2 Rotation;

// Output
out vec4 FragColor;

void main() {

    // Mix material color with textures colors
    vec4 texMixed = vec4(1);
    #if MAT_TEXTURES==1
        MIX_POINT_TEXTURE(0)
    #elif MAT_TEXTURES==2
        MIX_POINT_TEXTURE(0)
        MIX_POINT_TEXTURE(1)
    #elif MAT_TEXTURES==3
        MIX_POINT_TEXTURE(0)
        MIX_POINT_TEXTURE(1)
        MIX_POINT_TEXTURE(2)
    #endif

    // Generates final color
    FragColor = min(vec4(Color, MatOpacity) * texMixed, vec4(1));
}

`

const point_vertex_source = `#include <attributes>

// Model uniforms
uniform mat4 MVP;

// Material uniforms
#include <material>

// Outputs for fragment shader
out vec3 Color;
flat out mat2 Rotation;

void main() {

    // Rotation matrix for fragment shader
    float rotSin = sin(MatPointRotationZ);
    float rotCos = cos(MatPointRotationZ);
    Rotation = mat2(rotCos, rotSin, - rotSin, rotCos);

    // Sets the vertex position
    vec4 pos = MVP * vec4(VertexPosition, 1.0);
    gl_Position = pos;

    // Sets the size of the rasterized point decreasing with distance
    gl_PointSize = (1.0 - pos.z / pos.w) * MatPointSize;

    // Outputs color
    Color = MatEmissiveColor;
}

`

const sprite_fragment_source = `//
// Fragment shader for sprite
//

#include <material>

// Inputs from vertex shader
in vec3 Color;
in vec2 FragTexcoord;

// Output
out vec4 FragColor;

void main() {

    // Combine all texture colors and opacity
    vec4 texCombined = vec4(1);
#if MAT_TEXTURES>0
    for (int i = 0; i < {{.MatTexturesMax}}; i++) {
        vec4 texcolor = texture(MatTexture[i], FragTexcoord * MatTexRepeat(i) + MatTexOffset(i));
        if (i == 0) {
            texCombined = texcolor;
        } else {
            texCombined = mix(texCombined, texcolor, texcolor.a);
        }
    }
#endif

    // Combine material color with texture
    FragColor = min(vec4(Color, MatOpacity) * texCombined, vec4(1));
}

`

const sprite_vertex_source = `//
// Vertex shader for sprites
//

#include <attributes>

// Input uniforms
uniform mat4 MVP;

#include <material>

// Outputs for fragment shader
out vec3 Color;
out vec2 FragTexcoord;

void main() {

    // Applies transformation to vertex position
    gl_Position = MVP * vec4(VertexPosition, 1.0);

    // Outputs color
    Color = MatDiffuseColor;

    // Flips texture coordinate Y if requested.
    vec2 texcoord = VertexTexcoord;
#if MAT_TEXTURES>0
    if (MatTexFlipY[0]) {
        texcoord.y = 1 - texcoord.y;
    }
#endif
    FragTexcoord = texcoord;
}

`

const standard_fragment_source = `//
// Fragment Shader template
//
#include <material>

// Inputs from Vertex shader
in vec3 ColorFrontAmbdiff;
in vec3 ColorFrontSpec;
in vec3 ColorBackAmbdiff;
in vec3 ColorBackSpec;
in vec2 FragTexcoord;

// Output
out vec4 FragColor;


void main() {

    // Mix material color with textures colors
    vec4 texMixed = vec4(1);
    vec4 texColor;
    #if MAT_TEXTURES==1
        MIX_TEXTURE(0)
    #elif MAT_TEXTURES==2
        MIX_TEXTURE(0)
        MIX_TEXTURE(1)
    #elif MAT_TEXTURES==3
        MIX_TEXTURE(0)
        MIX_TEXTURE(1)
        MIX_TEXTURE(2)
    #endif

    vec4 colorAmbDiff;
    vec4 colorSpec;
    if (gl_FrontFacing) {
        colorAmbDiff = vec4(ColorFrontAmbdiff, MatOpacity);
        colorSpec = vec4(ColorFrontSpec, 0);
    } else {
        colorAmbDiff = vec4(ColorBackAmbdiff, MatOpacity);
        colorSpec = vec4(ColorBackSpec, 0);
    }
    FragColor = min(colorAmbDiff * texMixed + colorSpec, vec4(1));
}

`

const standard_vertex_source = `//
// Vertex shader standard
//
#include <attributes>

// Model uniforms
uniform mat4 ModelViewMatrix;
uniform mat3 NormalMatrix;
uniform mat4 MVP;

#include <lights>
#include <material>
#include <phong_model>


// Outputs for the fragment shader.
out vec3 ColorFrontAmbdiff;
out vec3 ColorFrontSpec;
out vec3 ColorBackAmbdiff;
out vec3 ColorBackSpec;
out vec2 FragTexcoord;

void main() {

    // Transform this vertex normal to camera coordinates.
    vec3 normal = normalize(NormalMatrix * VertexNormal);

    // Calculate this vertex position in camera coordinates
    vec4 position = ModelViewMatrix * vec4(VertexPosition, 1.0);

    // Calculate the direction vector from the vertex to the camera
    // The camera is at 0,0,0
    vec3 camDir = normalize(-position.xyz);

    // Calculates the vertex Ambient+Diffuse and Specular colors using the Phong model
    // for the front and back
    phongModel(position,  normal, camDir, MatAmbientColor, MatDiffuseColor, ColorFrontAmbdiff, ColorFrontSpec);
    phongModel(position, -normal, camDir, MatAmbientColor, MatDiffuseColor, ColorBackAmbdiff, ColorBackSpec);

    vec2 texcoord = VertexTexcoord;
#if MAT_TEXTURES > 0
    // Flips texture coordinate Y if requested.
    if (MatTexFlipY(0)) {
        texcoord.y = 1 - texcoord.y;
    }
#endif
    FragTexcoord = texcoord;

    gl_Position = MVP * vec4(VertexPosition, 1.0);
}

`

// Maps include name with its source code
var includeMap = map[string]string{

	"attributes":  include_attributes_source,
	"lights":      include_lights_source,
	"material":    include_material_source,
	"phong_model": include_phong_model_source,
}

// Maps shader name with its source code
var shaderMap = map[string]string{

	"basic_fragment":    basic_fragment_source,
	"basic_vertex":      basic_vertex_source,
	"panel_fragment":    panel_fragment_source,
	"panel_vertex":      panel_vertex_source,
	"phong_fragment":    phong_fragment_source,
	"phong_vertex":      phong_vertex_source,
	"physical_fragment": physical_fragment_source,
	"physical_vertex":   physical_vertex_source,
	"point_fragment":    point_fragment_source,
	"point_vertex":      point_vertex_source,
	"sprite_fragment":   sprite_fragment_source,
	"sprite_vertex":     sprite_vertex_source,
	"standard_fragment": standard_fragment_source,
	"standard_vertex":   standard_vertex_source,
}

// Maps program name with Proginfo struct with shaders names
var programMap = map[string]ProgramInfo{

	"basic":    {"basic_vertex", "basic_fragment", ""},
	"panel":    {"panel_vertex", "panel_fragment", ""},
	"phong":    {"phong_vertex", "phong_fragment", ""},
	"physical": {"physical_vertex", "physical_fragment", ""},
	"point":    {"point_vertex", "point_fragment", ""},
	"sprite":   {"sprite_vertex", "sprite_fragment", ""},
	"standard": {"standard_vertex", "standard_fragment", ""},
}
