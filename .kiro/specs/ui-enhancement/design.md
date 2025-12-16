# Design Document: UI Enhancement

## Overview

This design document outlines the approach for enhancing the visual design and user experience of the Easy Deploy terminal user interface. The enhancement will focus on creating a modern, cohesive, and beautiful terminal UI using the existing Bubble Tea and Lipgloss frameworks. The design maintains all current functionality while significantly improving the visual presentation through better color schemes, enhanced spacing, improved typography, and more polished visual elements.

The enhancement will be implemented by refactoring the existing style definitions and view rendering logic in `ui.go`, introducing new style components, and adding visual enhancements to all UI models (MainMenu, SSHList, SSHForm, SSHTest, ProjectList, ProjectForm, and Deploy).

## Architecture

The UI enhancement follows a style-first architecture where all visual styling is centralized and consistently applied across all views:

### Style System Architecture

```
Style Definitions (Global)
├── Color Palette (Primary, Secondary, Accent, Status colors)
├── Typography Styles (Titles, Body, Labels, Help text)
├── Component Styles (Borders, Containers, Forms, Lists)
└── Status Styles (Success, Error, Warning, Info, Loading)

View Models
├── MainMenuModel → Uses Menu Styles
├── SSHListModel → Uses List Styles
├── SSHFormModel → Uses Form Styles
├── SSHTestModel → Uses Status Styles
├── ProjectListModel → Uses List Styles
├── ProjectFormModel → Uses Form Styles
└── DeployModel → Uses Log/Status Styles
```

### Component Hierarchy

1. **Global Styles Layer**: Centralized style definitions using Lipgloss
2. **View Layer**: Individual model views that apply styles
3. **Rendering Layer**: View() methods that compose styled elements

## Components and Interfaces

### 1. Enhanced Style System

The style system will be expanded with the following components:

#### Color Palette
```go
// Primary colors for main UI elements
primaryColor := "#7C3AED"      // Vibrant purple
secondaryColor := "#10B981"    // Emerald green
accentColor := "#F59E0B"       // Amber
backgroundDark := "#1F2937"    // Dark gray
backgroundLight := "#374151"   // Medium gray

// Status colors
successColor := "#10B981"      // Green
errorColor := "#EF4444"        // Red
warningColor := "#F59E0B"      // Amber
infoColor := "#3B82F6"         // Blue
mutedColor := "#9CA3AF"        // Gray

// Text colors
textPrimary := "#F9FAFB"       // Almost white
textSecondary := "#D1D5DB"     // Light gray
textMuted := "#9CA3AF"         // Medium gray
```

#### Typography Styles
- **Title Style**: Large, bold, centered with gradient-like effect
- **Subtitle Style**: Medium, semi-bold for section headers
- **Body Style**: Regular weight for main content
- **Label Style**: Small, uppercase for form labels
- **Help Style**: Italic, muted for help text
- **Monospace Style**: For technical details (host, port, etc.)

#### Component Styles
- **Enhanced Border Style**: Multiple border variants (thick, thin, double, glow effect)
- **Card Style**: Container with shadow effect using box characters
- **Badge Style**: Small pill-shaped indicators for status
- **Divider Style**: Horizontal separators with decorative elements
- **Icon Style**: Unicode symbols for visual indicators

### 2. Enhanced View Components

#### MainMenuModel Enhancements
- Add welcome banner with ASCII art or decorative border
- Enhance menu items with icons/symbols
- Add subtle animations or visual separators
- Improve selection indicator with arrow or highlight box

#### SSHListModel Enhancements
- Add connection status indicators (colored dots)
- Enhance item cards with better spacing and borders
- Add visual grouping for related configurations
- Improve empty state with helpful message and styling

#### SSHFormModel Enhancements
- Add field type indicators (icons for password, key, etc.)
- Enhance cursor visibility with animated or styled indicator
- Add field validation visual feedback
- Improve field labels with better alignment and styling
- Add progress indicator showing form completion

#### SSHTestModel Enhancements
- Add animated spinner during connection test
- Enhance result display with icons and color coding
- Add connection details in a styled info box
- Improve error messages with helpful suggestions

#### ProjectListModel Enhancements
- Add project status badges
- Enhance server list display with chips/tags
- Add last deployment timestamp with relative time
- Improve visual hierarchy between project name and details

#### ProjectFormModel Enhancements
- Similar enhancements to SSHFormModel
- Add multi-line field support with better visualization
- Add field descriptions/hints below labels

#### DeployModel Enhancements
- Add deployment progress bar or step indicator
- Enhance log entries with timestamps and icons
- Add color-coded log levels (info, success, error)
- Add deployment summary at the end
- Improve real-time update visualization

### 3. Helper Functions

New helper functions to support enhanced UI:

```go
// renderIcon returns a styled icon/symbol
func renderIcon(iconType string, style lipgloss.Style) string

// renderBadge creates a pill-shaped badge
func renderBadge(text string, badgeType string) string

// renderDivider creates a decorative horizontal line
func renderDivider(width int, style lipgloss.Style) string

// renderProgressBar creates a progress indicator
func renderProgressBar(current, total int, width int) string

// renderSpinner returns animated spinner frame
func renderSpinner(frame int) string

// renderCard wraps content in a styled card
func renderCard(content string, title string) string

// renderKeyHelp formats keyboard shortcuts beautifully
func renderKeyHelp(keys map[string]string) string
```

## Data Models

No changes to existing data models are required. The enhancement focuses purely on visual presentation. However, we may add view-specific state for animations:

```go
// Add to models that need animation
type animationState struct {
    frame int
    ticker time.Ticker
}
```

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*


### Property Reflection

After reviewing all testable properties from the prework analysis, several properties can be consolidated:

- Properties 2.1 and 2.5 (consistent padding/margins and spacing ratios) can be combined into a single comprehensive property about consistent spacing
- Properties 3.1, 3.3, 3.4, and 3.5 (various border requirements) can be combined into a property about consistent border application
- Properties 4.1 and 4.5 (selected item highlighting) are redundant - one subsumes the other
- Properties 7.1, 7.3, and 7.5 (keyboard shortcut display) can be combined into one comprehensive property
- Properties 9.1 and 9.5 (list formatting and header styling) overlap and can be consolidated

The following properties represent the unique, non-redundant correctness requirements:

Property 1: Status message color differentiation
*For any* status message (success, error, warning, info), the rendered output should use the color designated for that message type
**Validates: Requirements 1.3**

Property 2: Consistent spacing across UI elements
*For any* UI element across all screens, the padding and margin values should follow consistent ratios and patterns
**Validates: Requirements 2.1, 2.5**

Property 3: Form field alignment
*For any* form, all form fields should be aligned consistently with uniform spacing between fields
**Validates: Requirements 2.3**

Property 4: Nested content indentation
*For any* nested content structure, the indentation level should increase proportionally with hierarchy depth
**Validates: Requirements 2.4**

Property 5: Consistent border styling
*For any* container, list, or form element, the border style (rounded, thickness, color) should be applied consistently based on element type
**Validates: Requirements 3.1, 3.3, 3.4, 3.5**

Property 6: Selected item highlighting
*For any* list with a selected item, the selected item should have distinct background color and visual indicator that differs from unselected items
**Validates: Requirements 4.1, 4.5**

Property 7: Active form field indication
*For any* form with an active field, the active field should display a cursor or highlight indicator
**Validates: Requirements 4.2**

Property 8: Title bold styling
*For any* title element, the rendered output should include bold styling
**Validates: Requirements 5.1**

Property 9: Description differentiation
*For any* description text, the styling (color or weight) should differ from primary text styling
**Validates: Requirements 5.2**

Property 10: Help text italic styling
*For any* help text element, the rendered output should include italic styling
**Validates: Requirements 5.3**

Property 11: Label consistency across forms
*For any* two forms, the label formatting (style, alignment, spacing) should be identical
**Validates: Requirements 5.4**

Property 12: Long text wrapping
*For any* text content exceeding the display width, the text should wrap to the next line with proper alignment maintained
**Validates: Requirements 5.5**

Property 13: Log entry type indicators
*For any* deployment log entry, the rendered output should include an icon or symbol corresponding to the entry type
**Validates: Requirements 6.1**

Property 14: Log entry alternating styles
*For any* sequence of multiple log entries, alternating visual styles or indentation should be applied for differentiation
**Validates: Requirements 6.5**

Property 15: Keyboard shortcut display consistency
*For any* screen, keyboard shortcuts should be displayed in a consistent location with symbols/formatting showing both key and action
**Validates: Requirements 7.1, 7.3, 7.5**

Property 16: Contextual help positioning
*For any* contextual help text, it should be positioned adjacent to or near the relevant content it describes
**Validates: Requirements 7.4**

Property 17: Form cursor position display
*For any* active form field, the current cursor position within the field should be visually indicated
**Validates: Requirements 8.1**

Property 18: Invalid field error styling
*For any* form field with validation errors, the field should be rendered with error styling
**Validates: Requirements 8.2**

Property 19: Multi-line field boundaries
*For any* multi-line form field, the field should have clear visual boundaries and adequate space
**Validates: Requirements 8.4**

Property 20: List item formatting consistency
*For any* list, all items should have identical formatting and alignment
**Validates: Requirements 9.1, 9.5**

Property 21: Supplementary information styling
*For any* list item with supplementary information, the supplementary content should use secondary styling distinct from primary content
**Validates: Requirements 9.2**

Property 22: Pagination indicator styling
*For any* paginated list, page indicators should be rendered with clear styling
**Validates: Requirements 9.4**

Property 23: Success message styling
*For any* successful action completion, a success message with positive/success styling should be displayed
**Validates: Requirements 10.1**

Property 24: Error message styling
*For any* error occurrence, an error message with error styling should be displayed
**Validates: Requirements 10.2**

Property 25: Processing indicator display
*For any* processing operation, a loading or progress indicator should be displayed
**Validates: Requirements 10.3**

Property 26: Input prompt styling
*For any* user input prompt, the prompt should be rendered with distinct styling
**Validates: Requirements 10.4**

Property 27: Status message visual distinction
*For any* status message, the styling should be visually distinct from non-status content
**Validates: Requirements 10.5**

## Error Handling

The UI enhancement focuses on visual presentation and does not introduce new error conditions. However, the enhanced UI will improve error presentation:

1. **Style Application Errors**: If a style fails to apply (e.g., invalid color), fall back to default styling
2. **Rendering Errors**: If enhanced rendering fails, fall back to basic text rendering
3. **Animation Errors**: If animation state fails, display static version

Error handling strategy:
- Graceful degradation: Always provide a working UI even if enhanced styles fail
- No crashes: Style failures should never crash the application
- Logging: Log style application failures for debugging

## Testing Strategy

### Unit Testing Approach

Unit tests will verify specific styling behaviors and edge cases:

1. **Style Definition Tests**
   - Verify color values are correctly defined
   - Test style composition (combining multiple styles)
   - Verify style inheritance and overrides

2. **Helper Function Tests**
   - Test icon rendering with various icon types
   - Test badge rendering with different badge types
   - Test divider rendering with various widths
   - Test progress bar calculation and rendering
   - Test card wrapping with different content sizes

3. **Edge Case Tests**
   - Empty list rendering (Requirement 9.3)
   - Deployment step completion display (Requirement 6.2)
   - Deployment step failure display (Requirement 6.3)
   - Deployment in progress display (Requirement 6.4)
   - Main menu header display (Requirement 3.2)

4. **Integration Tests**
   - Test complete view rendering for each model
   - Verify styles are applied correctly in context
   - Test view transitions maintain styling

### Property-Based Testing Approach

Property-based tests will verify universal styling properties across all inputs using the **gopter** library for Go. Each test will run a minimum of 100 iterations to ensure comprehensive coverage.

**Library Selection**: gopter (github.com/leanovate/gopter) - A mature property-based testing library for Go with good support for custom generators and properties.

**Test Configuration**: Each property test will be configured with:
```go
properties := gopter.NewProperties(parameters)
parameters.MinSuccessfulTests = 100
```

**Property Test Tagging**: Each property-based test will include a comment tag in this exact format:
```go
// **Feature: ui-enhancement, Property {number}: {property_text}**
```

**Property Test Implementation Requirements**:
1. Each correctness property listed above MUST be implemented by a SINGLE property-based test
2. Tests MUST generate random valid inputs (UI elements, content, states)
3. Tests MUST verify the property holds across all generated inputs
4. Tests MUST use the gopter library's property testing framework
5. Tests MUST be tagged with the property number and text from the design document

**Example Property Test Structure**:
```go
// **Feature: ui-enhancement, Property 1: Status message color differentiation**
func TestProperty_StatusMessageColors(t *testing.T) {
    properties := gopter.NewProperties(nil)
    properties.MinSuccessfulTests = 100
    
    properties.Property("status messages use designated colors", 
        prop.ForAll(
            func(msgType string, content string) bool {
                // Generate status message
                // Verify color matches message type
                return colorMatchesType(msgType, rendered)
            },
            gen.OneConstOf("success", "error", "warning", "info"),
            gen.AnyString(),
        ))
    
    properties.TestingRun(t)
}
```

### Testing Coverage

The dual testing approach ensures:
- **Unit tests** catch specific bugs in helper functions and edge cases
- **Property tests** verify styling rules hold universally across all UI states
- Together they provide confidence that the UI enhancement maintains consistency and correctness

## Implementation Notes

### Phase 1: Style System Enhancement
1. Define new color palette constants
2. Create enhanced typography styles
3. Implement helper functions for icons, badges, dividers
4. Create card and container styles

### Phase 2: View Enhancements
1. Update MainMenuModel view with enhanced styles
2. Update SSHListModel and SSHFormModel views
3. Update SSHTestModel with spinner and enhanced feedback
4. Update ProjectListModel and ProjectFormModel views
5. Update DeployModel with enhanced log visualization

### Phase 3: Polish and Refinement
1. Ensure consistent spacing across all views
2. Add animations where appropriate
3. Test on different terminal sizes
4. Optimize rendering performance

### Backward Compatibility

All enhancements maintain backward compatibility:
- No changes to data models or business logic
- No changes to keyboard shortcuts or navigation
- No changes to configuration file format
- Existing functionality remains unchanged

### Performance Considerations

- Style calculations are done once at initialization
- Avoid expensive string operations in hot paths
- Cache rendered components where possible
- Keep animation frame rates reasonable (avoid excessive redraws)

## Dependencies

- **Existing**: github.com/charmbracelet/bubbletea (TUI framework)
- **Existing**: github.com/charmbracelet/lipgloss (styling library)
- **New**: github.com/leanovate/gopter (property-based testing)

No additional runtime dependencies are required for the UI enhancements.
