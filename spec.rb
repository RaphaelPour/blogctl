require 'rspec'
require 'open3'
require 'tmpdir'

def blogctl(args, stdin: "")
  go_args = "ci-build/blogctl.test -test.run '^Test_'"\
    " -test.coverprofile=coverage/#{rand 0..10_000}.out"

  Open3.capture3({'TEST_ARGS' =>  args}, go_args, :stdin_data => stdin)
end

describe 'CLI' do

  context 'Trivial test cases' do
    it 'shows a version' do
      out, err, _ = blogctl("--version")
      expect(err).to eq ""
      expect(out).to match(/BuildVersion:/)
      expect(out).to match(/BuildDate:/)
    end

    it 'shows general help without any args' do
      out,err,_ = blogctl("")
      expect(err).to eq ""
      expect(out).to include 'Blogctl manages blog markdown-based posts'\
                             ' database-less and generates a static website on-demand'
      expect(out).to include 'Usage'
      expect(out).to include 'Available Commands'
      expect(out).to include 'Use "blogctl [command] --help" for more information about a command.'
    end

    it 'shows post help without any subcommand' do
      out,err,_ = blogctl("post")
      expect(err).to eq ""
      expect(out).to include("Manage posts")
      expect(out).to include 'Usage'
      expect(out).to include 'Available Commands'
      expect(out).to include 'Use "blogctl post [command] --help" for more information about a command.'

    end
  end

  context 'Basic test cases' do

    let (:blog_path) { "/tmp/blog#{rand 0..10_0000}" }
    
    after(:each) do
      FileUtils.rm_rf(blog_path)
    end

    it 'initializes a new blog' do
      _, err, _ = blogctl("init")
      expect(err).to eq ""
      expect(File.exist?("blog")).to be_truthy
    ensure
      FileUtils.rm_rf("blog")
    end

    it 'initializes a new blog with custom path' do
      _, err, _ = blogctl("init -p #{blog_path}")
      expect(err).to eq ""
      expect(File.exist?(blog_path)).to be_truthy
    end

    it 'fails to initializes a blog twice' do
      out,_,_ = blogctl("init -p #{blog_path}")
      expect(out).not_to match(/Blog environment already exists/)
      expect(File.exist?(blog_path)).to be_truthy
      
      out, _, _ = blogctl("init -p #{blog_path}")
      expect(out).to match(/Blog environment already exists/)
    end
  end

  context 'Post test cases' do

    let (:blog_path) { "/tmp/blog#{rand 0..10_0000}" }
    
    before(:each) do
      blogctl("init -p #{blog_path}")
    end

    after(:each) do
      FileUtils.rm_rf(blog_path)
    end

    it 'fails to add a post without title' do
      out, _, _ = blogctl("post add -p #{blog_path}")
      expect(out).to match(/Title missing/)
    end

    it 'adds a post' do
      out, _, _ = blogctl("post add --title 'test' -p #{blog_path}")
      expect(out).to include "#{blog_path}/test/content.md"
    end

    it 'updates a post' do
      blogctl("post add --title test -p #{blog_path}")

      appendix = "I use RSpec to test my blog"
      blogctl(
        "post update --slug test -p #{blog_path} -a",
        stdin: appendix
      )

      expect(File.read("#{blog_path}/test/content.md"))
        .to include(appendix)
    end

    it 'lists zero posts' do
      out, _, _ = blogctl("list -p #{blog_path}")
      expect(out).to match(/^Creation data\s*|\s*Title$/)
    end

    it 'lists one post' do
      blogctl("post add --title test -p #{blog_path}")
      out, _, _ = blogctl("list -p #{blog_path}")
      expect(out).to match(/^.*|\s*draft\s*|\s*test\s*$/)
    end

    it 'publishs a post' do
      blogctl("post add --title test -p #{blog_path}")
      out, _, _ = blogctl("list -p #{blog_path}")
      expect(out).to match(/^.*|\s*draft\s*|\s*test\s*$/)

      blogctl("post publish -p #{blog_path} --slug test")

      out, _, _ = blogctl("list -p #{blog_path}")
      expect(out).to match(/^.*|\s*public\s*|\s*test\s*$/)
    end

    it 'drafts a post again' do
      blogctl("post add --title test -p #{blog_path}")
      out, _, _ = blogctl("list -p #{blog_path}")
      blogctl("post publish -p #{blog_path} --slug test")

      out, _, _ = blogctl("list -p #{blog_path}")
      expect(out).to match(/^.*|\s*public\s*|\s*test\s*$/)

      blogctl("post draft -p #{blog_path} --slug test")

      out, _, _ = blogctl("list -p #{blog_path}")
      expect(out).to match(/^.*|\s*draft\s*|\s*test\s*$/)
    end
  end
end
