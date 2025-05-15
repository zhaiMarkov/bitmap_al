# README.md

## General Instructions


- Your code **MUST** be written in accordance with `gofumpt`. If not, you will be graded 0 automatically.
- Your program **MUST** be able to compile successfully.
- Your program **MUST NOT** exit unexpectedly (e.g., due to any panics like nil-pointer dereference, index out of range). If it does, you will receive a 0 during the defense.
- Only standard packages are allowed. If not, you will get 0 grade.
- The project MUST be compiled by the following command in the project's root directory:

```sh
$ go build -o bitmap .
```

- In case of an error, exit with a non-zero status and display a clear error message.

## Mandatory Features

## Header (Implemented by Maissaye)

The program must display bitmap properties with the command `header`. It should include:

- File type
- File size in bytes
- Header size
- DIB header size
- Width and height in pixels
- Pixel size in bits
- Image size in bytes

**Example command:**

```sh
$ ./bitmap header sample.bmp
BMP Header:
- FileType BM
- FileSizeInBytes 518456
- HeaderSize 54
DIB Header:
- DibHeaderSize 40
- WidthInPixels 480
- HeightInPixels 360
- PixelSizeInBits 24
- ImageSizeInBytes 518402
$
```

## core Package (implemented by Aomarbek)

The `core` package is designed to handle bitmap (BMP) image files, providing functionality to read, manipulate, and save bitmap images. It defines the structure of a BMP file, including its headers and pixel data.

### BitMap

The `BitMap` struct is the main representation of a bitmap image. It contains:

- `header`: An instance of `BMPHeader` that holds information about the BMP file.
- `infoHeader`: An instance of `DIBHeader` that contains detailed image information.
- `pixels`: A 2D slice of `Pixel` pointers representing the image's pixel data.

### BMPHeader

The `BMPHeader` struct contains metadata about the BMP file, such as:

- `FileType`: Indicates the file format (should be "BM").
- `FileSize`: The total size of the BMP file.
- `BitmapOffset`: The offset where the pixel data begins.

### DIBHeader

The `DIBHeader` struct contains details specific to the image, including:

- `Width` and `Height`: Dimensions of the image.
- `BitsPerPixel`: Number of bits used for each pixel (must be 24 for this implementation).
- `Compression`: Compression method (currently unhandled).
- `ImageSize`: Size of the image data.

### Pixel

The `Pixel` struct represents a single pixel's color, consisting of:

- `Blue`, `Green`, and `Red`: Color channels for the pixel.

## Methods

### NewBitMap

Creates and returns a new instance of `BitMap`.

### Read

Reads BMP file data from an `io.Reader` into the `BitMap` structure, handling both headers and pixel data. It accounts for padding required by BMP format.

### Save

Writes the current `BitMap` data back to an `io.Writer`, preserving the BMP file structure.

### Getters and Setters

The `BitMap` struct provides various getter and setter methods to access and modify header information, pixel data, and dimensions of the image.

### Read Methods for Headers and Pixel

The `BMPHeader`, `DIBHeader`, and `Pixel` structs each have their own `Read` method to facilitate reading their respective data from an `io.Reader`.

### Error Handling

The package includes error handling for invalid BMP formats, unsupported pixel formats, and issues during reading and writing operations. Errors are reported to standard error output, and the program exits on critical failures.

## mirror Package (Implemented by Missayev)

The `mirror` package provides functionality to mirror bitmap images either horizontally or vertically. It interacts with the core bitmap representation to perform the mirroring operations on the pixel data.

### Functions

### MirrorHorizontally

This function takes a 2D slice of `core.Pixel` pointers (representing the pixel data of a bitmap image) and returns a new 2D slice where the image is mirrored horizontally. 

- **Parameters**: 
  - `pixels`: A 2D slice containing the original pixel data.
  
- **Returns**: 
  - A new 2D slice of pixels that reflects the original image across the vertical axis.

### MirrorVertically

This function mirrors the pixel data of a bitmap image vertically. It also takes a 2D slice of `core.Pixel` pointers and returns a new mirrored version.

- **Parameters**: 
  - `pixels`: A 2D slice containing the original pixel data.
  
- **Returns**: 
  - A new 2D slice of pixels that reflects the original image across the horizontal axis.

### HandleMirror

The `HandleMirror` function processes the mirroring commands specified in the `config.MirrorFlag`. It modifies the pixel data of the provided `core.BitMap` based on the command.

- **Parameters**:
  - `bm`: A pointer to a `core.BitMap` instance representing the bitmap image.

- **Functionality**:
  - Checks if there are any mirroring commands in `config.MirrorFlag`.
  - Executes the corresponding mirroring function based on the command (either "horizontally" or "vertically").
  - If an invalid command is detected, it prints an error message and exits.
  - Updates the pixel data and dimensions of the bitmap accordingly.

### Error Handling

The package includes error handling for invalid mirroring commands, ensuring that the program exits gracefully if an unsupported option is provided.


## filter package (Implemented by Ykozhan)

The `filter` package is responsible for applying various image filters to bitmap images. It modifies the pixel data in a `core.BitMap` structure based on the specified filter commands.

### Filter Registry

A map called `filterRegistry` is used to associate filter names with their corresponding functions. Supported filters include:

- `blue`: Applies a blue filter.
- `red`: Applies a red filter.
- `green`: Applies a green filter.
- `grayscale`: Converts the image to grayscale.
- `negative`: Applies a negative effect to the image.
- `pixelate`: Applies a pixelation effect.
- `blur`: Applies a blur effect.

### HandleFilter

The `HandleFilter` function processes filter commands from the `config.FilterFlag`. It looks up the corresponding filter function in the `filterRegistry` and applies it to the provided bitmap.

- **Parameters**:
  - `b`: A pointer to a `core.BitMap` instance representing the bitmap image.

- **Functionality**:
  - Checks if there are any filter commands in `config.FilterFlag`.
  - If a valid filter is found, it calls the corresponding function; otherwise, it prints an error message and exits.

### Cycle Function

The `Cycle` helper function applies a transformation to each pixel in the bitmap using a provided function. It iterates through the pixel data and applies the specified operation.

### Filter Functions

Each filter function modifies the pixel data in specific ways:

- **ApplyRedFilter**: Sets the green and blue components of each pixel to zero.
- **ApplyGreenFilter**: Sets the red and blue components of each pixel to zero.
- **ApplyBlueFilter**: Sets the red and green components of each pixel to zero.
- **ApplyGrayscaleFilter**: Converts each pixel to grayscale using a weighted average based on human perception of color.
- **ApplyNegativeFilter**: Inverts the colors of each pixel by subtracting each color component from 255.
- **ApplyPixelateFilter**: Reduces detail by averaging colors in blocks of pixels and applying the average color to each pixel in that block. The block size increases with each application of the filter.
- **ApplyBlurFilter**: Blurs the image by averaging the color values of each pixel's neighbors in a defined range.

### Error Handling

The package includes error handling for unsupported filter commands, ensuring that the program exits gracefully if an invalid option is provided.

## rotate Package (Implemented by Maissyae)

The `rotate` package provides functionality to rotate bitmap images by specified angles. It modifies the pixel data in a `core.BitMap` structure according to the rotation commands.

### Global Variables

- `globalHeight` and `globalWidth`: These variables hold the current dimensions of the bitmap image during rotation operations.

### Rotation Map

A map called `rotationMap` associates rotation commands (as strings) with the number of 90-degree clockwise rotations needed. Supported rotations include:

- `right`, `90`, `-270`: 3 clockwise rotations (270 degrees counter-clockwise).
- `left`, `270`, `-90`: 1 clockwise rotation (90 degrees counter-clockwise).
- `180`, `-180`: 2 clockwise rotations (180 degrees).

### RotateBMP Function

The `RotateBMP` function takes a 2D slice of `core.Pixel` pointers representing the image's pixel data and rotates it 90 degrees clockwise.

- **Parameters**:
  - `pixels`: A 2D slice of pixels to rotate.

- **Returns**:
  - A new 2D slice of rotated pixels.

### rotateImage Function

The `rotateImage` function repeatedly calls `RotateBMP` for a specified number of rotations.

- **Parameters**:
  - `pixels`: The original pixel data.
  - `rotations`: The number of 90-degree clockwise rotations to apply.

- **Returns**:
  - The rotated pixel data.

### HandleRotate Function

The `HandleRotate` function processes rotation commands from the `config.RotateFlag`. It determines the appropriate rotation based on the command and applies the rotation to the bitmap.

- **Parameters**:
  - `b`: A pointer to a `core.BitMap` instance representing the bitmap image.

- **Functionality**:
  - Checks if there are any rotation commands in `config.RotateFlag`.
  - Looks up the number of rotations in the `rotationMap`. If the command is valid, it applies the corresponding rotations.
  - Updates the bitmap's pixel data and dimensions accordingly.
  - Exits with an error message if the rotation command is invalid.

### Error Handling

The package includes error handling for invalid rotation commands, ensuring that the program exits gracefully if an unsupported rotation is specified.


## crop package (Implemented by Mduisen)

The `crop` package provides functionality to crop bitmap images by defining a specific rectangular region. It modifies the pixel data in a `core.BitMap` structure based on the specified crop parameters.

### CropValues Struct

A structure called `CropValues` stores the parameters needed for cropping:

- `OffSetX`: The x-coordinate to start the crop from.
- `OffSetY`: The y-coordinate to start the crop from.
- `Width`: The width of the cropped area.
- `Height`: The height of the cropped area.

### HandleCrop Function

The `HandleCrop` function processes crop commands from the `config.CropFlag`. It extracts and validates cropping parameters, then modifies the bitmap accordingly.

- **Parameters**:
  - `b`: A pointer to a `core.BitMap` instance representing the bitmap image.

- **Functionality**:
  - Checks if there are crop commands in `config.CropFlag`.
  - Parses the crop command to extract the crop values.
  - Validates the crop values to ensure they are within the image dimensions.
  - Calls `cropImage` to modify the pixel data based on the specified crop values.
  - Updates the bitmap's dimensions and sizes.

### parseFlags Function

The `parseFlags` function splits the crop flag string into individual values and converts them into integers.

- **Parameters**:
  - `flag`: The crop command string.

- **Returns**:
  - A slice of strings containing the parsed values or an error if the parsing fails.

### handleTwoValues and handleFourValues Functions

These helper functions process the parsed crop values:

- **handleTwoValues**: Extracts the x and y offsets and calculates the width and height based on the image dimensions.
- **handleFourValues**: Extracts all four crop parameters (offsets and dimensions).

### validate Function

The `validate` function checks if the crop parameters are valid, ensuring that they are positive and do not exceed the image dimensions.

### convertToInt Function

Converts a slice of strings to a slice of integers for further processing.

### cropImage Function

The `cropImage` function performs the actual cropping of the image based on the provided crop parameters.

- **Parameters**:
  - `b`: A pointer to a `core.BitMap`.
  - `cropValues`: The crop parameters defined in `CropValues`.

- **Functionality**:
  - Creates a new 2D slice for cropped pixels.
  - Iterates through the original pixel data to fill in the cropped pixel data.
  - Updates the bitmap's dimensions and calculates the new image size.
  - Adjusts the bitmap's file size to reflect the cropped image.
  - Sets the new pixel data in the bitmap.

### Error Handling

The package includes error handling for invalid crop commands and values, ensuring that the program exits gracefully if the provided parameters are incorrect.

# config Package (implemented by Aomarbek)

The `config` package is responsible for handling command-line flags and arguments for a bitmap image processing application. It defines various commands, manages flag parsing, and validates input to ensure proper usage.

### Help Text

The package contains help text strings for guiding users on how to use the application:

- `helpText`: General usage information for the application.
- `headerHelpText`: Usage information specifically for the `header` command.
- `applyHelpText`: Usage information for the `apply` command, including options for mirroring, filtering, rotating, and cropping images.

### Command Handling

The package defines a map `m` that associates commands with their respective handler functions:

- `"header"`: Calls `handleHeader`.
- `"apply"`: Calls `handleApply`.

### Flag Sets

Two `flag.FlagSet` instances are created for handling command-specific flags:

- `HeaderCmd`: For the `header` command.
- `ApplyCmd`: For the `apply` command, which includes several image processing options.

### Flag Variables

The package defines variables to hold flag values for various processing options:

- `MirrorFlag`: A slice of strings for mirror operations.
- `FilterFlag`: A slice of strings for filter operations.
- `RotateFlag`: A slice of strings for rotation operations.
- `CropFlag`: A slice of strings for crop operations.
- `SourceFileName`: The name of the source bitmap file.
- `OutputFileName`: The name of the output bitmap file.
- `OrderedFlags`: A slice to maintain the order of flags passed.

### InitFlags Function

The `InitFlags` function initializes flag parsing:

- It checks if there are enough command-line arguments and displays usage information if not.
- It looks up the command in the map and calls the corresponding handler.

### handleHeader and handleApply Functions

These functions manage the specifics of the `header` and `apply` commands:

- **handleHeader**:
  - Sets up the `HeaderCmd` flag set.
  - Parses flags and validates input.
  - Sets the `SourceFileName` for the bitmap file to read.

- **handleApply**:
  - Sets up the `ApplyCmd` flag set with various image processing options.
  - Parses flags and validates input.
  - Sets the `SourceFileName` and `OutputFileName` for the bitmap files.

### parseFlags Function

The `parseFlags` function handles the parsing of command-specific flags:

- It checks if there are enough arguments for the command and exits if not.
- It uses the `Parse` method of the `FlagSet` to parse the command-line arguments.

### parseOrderedFlags Function

This function extracts and stores flags that are prefixed with `--` or `-` into the `OrderedFlags` slice, preserving the order they were provided.

### Validation Functions

Two validation functions ensure the integrity of command inputs:

- **validateHeader**:
  - Validates the arguments for the `header` command.
  - Checks for correct argument count and file format.

- **validateApply**:
  - Validates the arguments for the `apply` command.
  - Ensures there are exactly two file arguments and checks their formats.

### hasFlags Function

The `hasFlags` function checks if any argument in a given slice starts with a `-`, indicating the presence of flags.

## Error Handling

The package includes robust error handling that prints error messages and exits the application if user input does not meet the expected criteria.

# help Package (implemented by Aomarbek)

The `help` package is responsible for providing usage instructions and descriptions for the bitmap image processing application. It defines help texts for the overall program and specific commands to guide users on how to use the application effectively.

### Help Texts

1. **General Help Text (`helpText`)**:
   - Provides an overview of the application usage.
   - Lists available commands, including:
     - `header`: Prints bitmap file header information.
     - `apply`: Applies processing to the image and saves it to a file.

   ```plaintext
   Usage:
     bitmap <command> [arguments]

   The commands are:
     header    prints bitmap file header information
     apply     applies processing to the image and saves it to the file

2. **Header Help Text (headerHelpText)**:

    -Provides usage information specifically for the header command.

```plaintext
Usage:
  bitmap header <filename>
```

3. **Apply Help Text (applyHelpText)**:

    -Provides usage information for the apply command, including options for mirroring, filtering, rotating, and cropping images.
```plaintext
Usage:
  bitmap apply <command> <source file> <output file>

Available commands:
  mirror     mirrors the image (options: horizontally, vertically)
  filter     applies a color filter (options: red, green, blue, grayscale, negative, pixelate, blur)
  rotate     rotates the image (options: right, left, 180)
  crop       crops the image (specify dimensions)
```
### Functionality

The help package ensures that users can access detailed command usage information, allowing them to effectively utilize the bitmap image processing application.